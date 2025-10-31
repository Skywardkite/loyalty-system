package service

import (
	"context"

	"github.com/Skywardkite/loyalty-system/internal/handler/dto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrBadRequest   = errors.New("Некорректный запрос")
	ErrUnauthorized = errors.New("Нет доступа")
)

func (s *Service) Login(ctx context.Context, req *dto.User) (int64, error) {
	if req.Login == "" || req.Password == "" {
		return 0, ErrBadRequest
	}

	userID, hashedPassword, err := s.store.GetUserByLogin(ctx, req.Login)
	if err != nil {
		s.logger.Errorw("failed to get user by login", "error", err)
		return 0, err
	}

	// Сравниваем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return 0, ErrUnauthorized
	}

	return userID, nil
}
