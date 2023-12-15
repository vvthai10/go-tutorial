package entry

import (
	"context"

	"github.com/vvthai10/transaction-mongodb/entities"
)

type EntryService struct {
	repo IEntryRepository
}

func NewEntryService(repo IEntryRepository) IEntryService {
	return &EntryService{
		repo: repo,
	}
}

func (s *EntryService) CreateEntry(ctx context.Context, arg CreateEntryParams) (entities.Entry, error) {
	return s.repo.CreateEntry(ctx, arg)
}

func (s *EntryService) GetEntry(ctx context.Context, id string) (entities.Entry, error) {
	return s.repo.GetEntry(ctx, id)
}
