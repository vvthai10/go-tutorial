package transfer

import (
	"context"

	"github.com/vvthai10/transaction-mongodb/entities"
)

type TransferService struct {
	repo ITransferRepository
}

func NewTransferService(repo ITransferRepository) ITransferService {
	return &TransferService{
		repo: repo,
	}
}

func (s *TransferService) CreateTransfer(ctx context.Context, arg CreateTransferParams) (CreateTransferTxResult, error) {
	return s.repo.CreateTransferTx(ctx, arg)
}

func (s *TransferService) GetTransfer(ctx context.Context, id string) (entities.Transfer, error) {
	return s.repo.GetTransfer(ctx, id)
}
