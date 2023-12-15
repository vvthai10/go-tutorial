package account_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vvthai10/transaction-mongodb/bootstrap"
	"github.com/vvthai10/transaction-mongodb/config"
	"github.com/vvthai10/transaction-mongodb/entities"
	"github.com/vvthai10/transaction-mongodb/repository"
	"github.com/vvthai10/transaction-mongodb/service/account"
	"github.com/vvthai10/transaction-mongodb/util"
)

var accountService account.IAccountService

func TestMain(m *testing.M) {
	env := &config.Env{
		DBUri: "mongodb+srv://vvthai1410:vvthai1410@cluster0.bunlj.mongodb.net/?retryWrites=true&w=majority",
	}
	db := bootstrap.NewMongoDB(env)
	accountRepo := repository.NewAccountRepository(db)
	accountService = account.NewAccountService(accountRepo)
	os.Exit(m.Run())
}

func createRandomAccount(t *testing.T) entities.Account {
	arg := account.CreateAccountParams {
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
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

func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	// create account
	account1 := createRandomAccount(t)
	account2, err := accountService.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
}