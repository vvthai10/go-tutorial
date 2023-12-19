package user

import (
	"context"

	"github.com/vvthai10/demo-graphql/entities"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}

type IUserService interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}

type CreateUserParams struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
