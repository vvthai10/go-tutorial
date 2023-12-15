package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vvthai10/transaction-mongodb/entities"
	"github.com/vvthai10/transaction-mongodb/service/account"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountRepository struct {
	db *mongo.Database
}

func NewAccountRepository(db *mongo.Database) account.IAccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, arg account.CreateAccountParams) (entities.Account, error) {
	account := entities.Account{
		Owner:    arg.Owner,
		Balance:  arg.Balance,
		Currency: arg.Currency,
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	inserted, err := r.db.Collection("accounts").InsertOne(ctx, account)
	if err != nil {
		return entities.Account{}, err
	}
	account.ID = inserted.InsertedID.(primitive.ObjectID).Hex()

	return account, nil
}

func (r *AccountRepository) GetAccount(ctx context.Context, id string) (entities.Account, error) {
	account := entities.Account{}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	_id, err := primitive.ObjectIDFromHex(id)
	err = r.db.Collection("accounts").FindOne(ctx, bson.M{"_id": _id}).Decode(&account)
	if err != nil {
		fmt.Println(err)
		return entities.Account{}, err
	}

	return account, nil
}

func (r *AccountRepository) AddAccountBalance(ctx context.Context, arg account.AddAccountBalanceParams) (entities.Account, error) {
	_id, _ := primitive.ObjectIDFromHex(arg.ID)
	filter := bson.M{"_id": _id}
	update := bson.M{"$inc": bson.M{"balance": arg.Amount}}
	account := entities.Account{}
	err := r.db.Collection("accounts").FindOneAndUpdate(ctx, filter, update).Decode(&account)
	txName := ctx.Value(TxKey)
	fmt.Println(txName, " account info ", account.Balance, account.Owner)
	if err != nil {
		fmt.Println(err)
		return entities.Account{}, err
	}

	return account, nil
}
