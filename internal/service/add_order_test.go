package service

import (
	"context"
	"testing"

	"github.com/Skywardkite/loyalty-system/internal/repository"
	repErr "github.com/Skywardkite/loyalty-system/internal/repository/error"
	"github.com/Skywardkite/loyalty-system/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestService_AddOrder(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		store  repository.Storage
	}

	type args struct {
		number string
		userID int64
	}

	userID := int64(1)
	number := "2377225624"
	repoErr := errors.New("user not found")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetOrderUserByNumber(mock.Anything, number).Return(0, repErr.ErrOrderNotExsit)
					repoMock.EXPECT().AddOrder(mock.Anything, userID, number).Return(nil)
					return repoMock
				}(),
			},
			args: args{
				number: number,
				userID: userID,
			},
			wantErr: nil,
		},
		{
			name: "bad request",
			args: args{
				number: "",
				userID: userID,
			},
			wantErr: ErrBadRequest,
		},
		{
			name: "invalid order number",
			args: args{
				number: "123",
				userID: userID,
			},
			wantErr: ErrInvalidOrderNumber,
		},
		{
			name: "err - get order",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetOrderUserByNumber(mock.Anything, number).Return(0, repoErr)
					return repoMock
				}(),
			},
			args: args{
				number: number,
				userID: userID,
			},
			wantErr: repoErr,
		},
		{
			name: "err - already uploaded by this user",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetOrderUserByNumber(mock.Anything, number).Return(userID, nil)
					return repoMock
				}(),
			},
			args: args{
				number: number,
				userID: userID,
			},
			wantErr: ErrOrderUploadedThisUser,
		},
		{
			name: "err - already uploaded by another user",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetOrderUserByNumber(mock.Anything, number).Return(12, nil)
					return repoMock
				}(),
			},
			args: args{
				number: number,
				userID: userID,
			},
			wantErr: ErrOrderUploadedAnotherUser,
		},
		{
			name: "err add order",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().GetOrderUserByNumber(mock.Anything, number).Return(0, repErr.ErrOrderNotExsit)
					repoMock.EXPECT().AddOrder(mock.Anything, userID, number).Return(repoErr)
					return repoMock
				}(),
			},
			args: args{
				number: number,
				userID: userID,
			},
			wantErr: repoErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				logger: tt.fields.logger,
				store:  tt.fields.store,
			}

			err := s.AddOrder(context.Background(), tt.args.number, tt.args.userID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
