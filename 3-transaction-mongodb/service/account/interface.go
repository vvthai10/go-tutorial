package account

import (
	"context"

	"github.com/vvthai10/transaction-mongodb/entities"
)

type IAccountRepository interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (entities.Account, error)
	GetAccount(ctx context.Context, id string) (entities.Account, error)
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (entities.Account, error)
}

type IAccountService interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (entities.Account, error)
	GetAccount(ctx context.Context, id string) (entities.Account, error)
}

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type AddAccountBalanceParams struct {
	ID     string `bson:"_id"`
	Amount int64  `bson:"balance"`
}
