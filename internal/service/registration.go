package service

import (
	"context"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlredyExists = errors.New("Пользователь уже существует в системе")
)

func (s *Service) Registration(ctx context.Context, req *dto.User) (int64, error) {
	if req.Login == "" || req.Password == "" {
		return 0, ErrBadRequest
	}

	exists, err := s.store.UserExists(ctx, req.Login)
	if err != nil {
		s.logger.Errorw("failed to check user existence", "error", err)
		return 0, err
	}
	if exists {
		return 0, ErrUserAlredyExists
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorw("failed to hash password", "error", err)
		return 0, err
	}

	// Создаём пользователя
	userID, err := s.store.CreateUser(ctx, req.Login, string(hashedPassword))
	if err != nil {
		s.logger.Errorw("failed to create user", "error", err)
		return 0, err
	}

	return userID, nil
}
