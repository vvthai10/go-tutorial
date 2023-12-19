package postgres

import (
	"context"

	"github.com/vvthai10/demo-graphql/entities"
	"github.com/vvthai10/demo-graphql/service/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, arg user.CreateUserParams) (*entities.User, error) {
	user := entities.User{
		FullName: arg.FullName,
		Email:    arg.Email,
		Password: arg.Password,
	}
	err := r.db.Create(&user).Error
	if err != nil {
		return &entities.User{}, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := entities.User{}
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return &entities.User{}, err
	}

	return &user, nil
}
