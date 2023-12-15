package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vvthai10/transaction-mongodb/entities"
	"github.com/vvthai10/transaction-mongodb/service/entry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EntryRepository struct {
	db *mongo.Database
}

func NewEntryRepository(db *mongo.Database) entry.IEntryRepository {
	return &EntryRepository{
		db: db,
	}
}

func (r *EntryRepository) CreateEntry(ctx context.Context, arg entry.CreateEntryParams) (entities.Entry, error) {
	entry := entities.Entry{
		AccountID: arg.AccountID,
		Amount:        arg.Amount,
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	inserted, err := r.db.Collection("entries").InsertOne(ctx, entry)
	if err != nil {
		return entities.Entry{}, err
	}
	entry.ID = inserted.InsertedID.(primitive.ObjectID).Hex()

	return entry, nil
}

func (r *EntryRepository) GetEntry(ctx context.Context, id string) (entities.Entry, error) {
	entry := entities.Entry{}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	_id, _ := primitive.ObjectIDFromHex(id)
	err := r.db.Collection("entries").FindOne(ctx, bson.M{"_id": _id}).Decode(&entry)
	if err != nil {
		fmt.Println(err)
		return entities.Entry{}, err
	}

	return entry, nil
}