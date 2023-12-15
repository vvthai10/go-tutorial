package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vvthai10/transaction-mongodb/entities"
	"github.com/vvthai10/transaction-mongodb/service/account"
	"github.com/vvthai10/transaction-mongodb/service/entry"
	"github.com/vvthai10/transaction-mongodb/service/transfer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransferRepository struct {
	db *mongo.Database
}

func NewTransferRepository(db *mongo.Database) transfer.ITransferRepository {
	return &TransferRepository{
		db: db,
	}
}

func (r *TransferRepository) ExecuteTransaction(ctx context.Context, fn func(account.IAccountRepository, entry.IEntryRepository, transfer.ITransferRepository) error) error {
	var session mongo.Session
	var err error

	if session, err = r.db.Client().StartSession(); err != nil {
		fmt.Println(err)
	}
	if err = session.StartTransaction(); err != nil {
		fmt.Println(err)
	}
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		
		accountRepo := NewAccountRepository(sc.Client().Database("transaction-db"))
		entryRepo := NewEntryRepository(sc.Client().Database("transaction-db"))
		transferRepo := NewTransferRepository(sc.Client().Database("transaction-db"))
		
		err = fn(accountRepo, entryRepo, transferRepo)

		if err != nil {
			if err = session.AbortTransaction(sc); err != nil {
				fmt.Println(err)
			}
			return err
		}
		
		if err = session.CommitTransaction(sc); err != nil {
			fmt.Println(err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	session.EndSession(ctx)

	return nil
}

var TxKey = struct{}{}

func (r *TransferRepository) CreateTransferTx(ctx context.Context, arg transfer.CreateTransferParams) (transfer.CreateTransferTxResult, error) {
	var result transfer.CreateTransferTxResult

	err := r.ExecuteTransaction(ctx, func(accountRepo account.IAccountRepository, entryRepo entry.IEntryRepository, transferRepo transfer.ITransferRepository) error {
		var err error

		txName := ctx.Value(TxKey)
		fmt.Println(txName, " create transfer")

		result.Transfer, err = transferRepo.CreateTransfer(ctx, transfer.CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, " create entry 1")
		result.FromEntry, err = entryRepo.CreateEntry(ctx, entry.CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println(txName, " create entry 2")
		result.ToEntry, err = entryRepo.CreateEntry(ctx, entry.CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			fmt.Println(txName, " add money 1")
			result.FromAccount, result.ToAccount, _ = addMoney(ctx, accountRepo, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			fmt.Println(txName, " add money 2")
			result.ToAccount, result.FromAccount, _ = addMoney(ctx, accountRepo, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})

	return result, err
}

func (r *TransferRepository) CreateTransfer(ctx context.Context, arg transfer.CreateTransferParams) (entities.Transfer, error) {
	transfer := entities.Transfer{
		FromAccountID: arg.FromAccountID,
		ToAccountID:   arg.ToAccountID,
		Amount:        arg.Amount,
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	inserted, err := r.db.Collection("transfers").InsertOne(ctx, transfer)
	if err != nil {
		return entities.Transfer{}, err
	}
	transfer.ID = inserted.InsertedID.(primitive.ObjectID).Hex()
	return transfer, nil
}

func (r *TransferRepository) GetTransfer(ctx context.Context, id string) (entities.Transfer, error) {
	transfer := entities.Transfer{}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	_id, _ := primitive.ObjectIDFromHex(id)
	err := r.db.Collection("transfers").FindOne(ctx, bson.M{"_id": _id}).Decode(&transfer)
	if err != nil {
		fmt.Println(err)
		return entities.Transfer{}, err
	}

	return transfer, nil
}

func addMoney(ctx context.Context, accountRepo account.IAccountRepository, accountID1 string, amount1 int64, accountID2 string, amount2 int64) (account1, account2 entities.Account, err error) {
	account1, err = accountRepo.AddAccountBalance(ctx, account.AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = accountRepo.AddAccountBalance(ctx, account.AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
