package transfer_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vvthai10/transaction-mongodb/bootstrap"
	"github.com/vvthai10/transaction-mongodb/config"
	"github.com/vvthai10/transaction-mongodb/entities"
	"github.com/vvthai10/transaction-mongodb/repository"
	"github.com/vvthai10/transaction-mongodb/service/account"
	"github.com/vvthai10/transaction-mongodb/service/entry"
	"github.com/vvthai10/transaction-mongodb/service/transfer"
	"github.com/vvthai10/transaction-mongodb/util"
)

var transferService transfer.ITransferService
var entryService entry.IEntryService
var accountService account.IAccountService

func TestMain(m *testing.M) {
	env := &config.Env{
		DBUri: "mongodb+srv://vvthai1410:vvthai1410@cluster0.bunlj.mongodb.net/?retryWrites=true&w=majority",
	}
	db := bootstrap.NewMongoDB(env)

	transferRepo := repository.NewTransferRepository(db)
	transferService = transfer.NewTransferService(transferRepo)

	entryRepo := repository.NewEntryRepository(db)
	entryService = entry.NewEntryService(entryRepo)

	accountRepo := repository.NewAccountRepository(db)
	accountService = account.NewAccountService(accountRepo)
	os.Exit(m.Run())
}

func createRandomAccount(t *testing.T) entities.Account {
	arg := account.CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := accountService.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	return account
}

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("[before] ", account1.Balance, account2.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan transfer.CreateTransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), repository.TxKey, txName)
			result, err := transferService.CreateTransfer(ctx, transfer.CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)

		_, err = transferService.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)

		_, err = entryService.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)

		_, err = entryService.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check account's balance
		fmt.Println("[tx] ", fromAccount.Balance, toAccount.Balance, " check ", account1.Balance, account2.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 >= 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount1, err := accountService.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := accountService.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("[after] ", account1.Balance, account2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}
