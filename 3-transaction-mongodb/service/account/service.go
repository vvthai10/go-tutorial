package account

import (
	"context"

	"github.com/vvthai10/transaction-mongodb/entities"
)

type AccountService struct {
	repo IAccountRepository
}

func NewAccountService(repo IAccountRepository) IAccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, arg CreateAccountParams) (entities.Account, error) {
	return s.repo.CreateAccount(ctx, arg)
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (entities.Account, error) {
	return s.repo.GetAccount(ctx, id)
}
