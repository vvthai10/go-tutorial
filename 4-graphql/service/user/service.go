package user

import (
	"context"

	"github.com/vvthai10/demo-graphql/entities"
)

type UserService struct {
	repo IUserRepository
}

func NewService(r IUserRepository) IUserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) CreateUser(ctx context.Context, arg CreateUserParams) (*entities.User, error) {
	return s.repo.CreateUser(ctx, arg)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
