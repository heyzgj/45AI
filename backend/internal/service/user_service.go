package service

import (
	"context"

	"github.com/45ai/backend/internal/model"
	"github.com/45ai/backend/internal/repository"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
} 