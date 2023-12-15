package entry

import (
	"context"

	"github.com/vvthai10/transaction-mongodb/entities"
)

type IEntryRepository interface {
	CreateEntry(ctx context.Context, arg CreateEntryParams) (entities.Entry, error)
	GetEntry(ctx context.Context, id string) (entities.Entry, error)
}

type IEntryService interface {
	CreateEntry(ctx context.Context, arg CreateEntryParams) (entities.Entry, error)
	GetEntry(ctx context.Context, id string) (entities.Entry, error)
}

type CreateEntryParams struct {
	AccountID 	string 	`json:"accountId"`
	Amount 			int64 	`json:"amount"`
}