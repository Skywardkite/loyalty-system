package service

import (
	"context"
	"testing"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"github.com/Skywardkite/loyalty-system/internal/repository"
	"github.com/Skywardkite/loyalty-system/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go.uber.org/zap"
)

func TestService_WriteOffPoints(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		store  repository.Storage
	}

	type args struct {
		req *dto.ParamsWithdraw
	}

	userID := int64(1)
	number := "2377225624"

	req := dto.ParamsWithdraw{
		Order: number,
		Sum:   78.90,
	}

	repErr := errors.New("repository error")

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
					repoMock.EXPECT().AddWithdrawal(mock.Anything, userID, int64(7890), req.Order).Return(nil)
					return repoMock
				}(),
			},
			args:    args{req: &req},
			wantErr: nil,
		},
		{
			name:    "bad request - empty order",
			args:    args{req: &dto.ParamsWithdraw{Order: "", Sum: 78.90}},
			wantErr: ErrUnprocessableEntity,
		},
		{
			name:    "bad request - empty sum",
			args:    args{req: &dto.ParamsWithdraw{Order: number, Sum: -78.90}},
			wantErr: ErrUnprocessableEntity,
		},
		{
			name: "error - user exists",
			fields: fields{
				logger: zap.NewNop().Sugar(),
				store: func() repository.Storage {
					repoMock := mocks.NewMockStorage(t)
					repoMock.EXPECT().AddWithdrawal(mock.Anything, userID, int64(7890), req.Order).Return(repErr)
					return repoMock
				}(),
			},
			args:    args{req: &req},
			wantErr: repErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				logger: tt.fields.logger,
				store:  tt.fields.store,
			}

			err := s.WriteOffPoints(context.Background(), tt.args.req, userID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
