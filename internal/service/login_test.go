package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"github.com/Skywardkite/loyalty-system/internal/repository"
	"github.com/Skywardkite/loyalty-system/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"go.uber.org/zap"
)

func TestService_Login(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		store  repository.Storage
	}

	type args struct {
		req *dto.User
	}

	req := dto.User{
		Login:    "userLogin",
		Password: "123456",
	}

	userID := int64(1)
	hashed, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	repErr := errors.New("user not found")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetUserByLogin(mock.Anything, req.Login).Return(userID, string(hashed), nil)
					return repoMock
				}(),
			},
			args:    args{req: &req},
			want:    userID,
			wantErr: nil,
		},
		{
			name:    "bad request - empty password",
			args:    args{req: &dto.User{Login: "userLogin", Password: ""}},
			want:    0,
			wantErr: ErrBadRequest,
		},
		{
			name:    "bad request - empty login",
			args:    args{req: &dto.User{Login: "", Password: "123456"}},
			want:    0,
			wantErr: ErrBadRequest,
		},
		{
			name: "error - wrong password",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetUserByLogin(mock.Anything, req.Login).Return(userID, string(hashed), nil)
					return repoMock
				}(),
			},
			args:    args{req: &dto.User{Login: "userLogin", Password: "wrongPassword"}},
			want:    0,
			wantErr: ErrUnauthorized,
		},
		{
			name: "error - user not found",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetUserByLogin(mock.Anything, req.Login).Return(0, "", repErr)
					return repoMock
				}(),
			},
			args:    args{req: &req},
			want:    0,
			wantErr: repErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				logger: tt.fields.logger,
				store:  tt.fields.store,
			}

			got, err := s.Login(context.Background(), tt.args.req)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.True(t, reflect.DeepEqual(got, tt.want))
		})
	}
}
