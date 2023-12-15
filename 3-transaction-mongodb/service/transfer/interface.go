package transfer

import (
	"context"

	"github.com/vvthai10/transaction-mongodb/entities"
	"github.com/vvthai10/transaction-mongodb/service/account"
	"github.com/vvthai10/transaction-mongodb/service/entry"
)

type ITransferRepository interface {
	ExecuteTransaction(ctx context.Context, fn func(account.IAccountRepository, entry.IEntryRepository, ITransferRepository) error) error
	CreateTransferTx(ctx context.Context, arg CreateTransferParams) (CreateTransferTxResult, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (entities.Transfer, error)
	GetTransfer(ctx context.Context, id string) (entities.Transfer, error)
}

type ITransferService interface {
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (CreateTransferTxResult, error)
	GetTransfer(ctx context.Context, id string) (entities.Transfer, error)
}

type CreateTransferParams struct {
	FromAccountID string `json:"fromAccountId"`
	ToAccountID   string `json:"toAccountId"`
	Amount        int64  `json:"amount"`
}

type CreateTransferTxResult struct {
	Transfer    entities.Transfer `json:"transfer"`
	FromAccount entities.Account  `json:"fromAccount"`
	ToAccount   entities.Account  `json:"toAccount"`
	FromEntry   entities.Entry    `json:"fromEntry"`
	ToEntry     entities.Entry    `json:"toEntry"`
}
