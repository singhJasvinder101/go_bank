package db

import (
	"context"
	"testing"
	"time"

	"github.com/singhJasvinder101/go_bank/utils"
	"github.com/stretchr/testify/require"
)

// transactions concurrency testings

func createRandomAccount(t *testing.T) Account{
	arg := CreateAccountParams{
		Owner: utils.RandomOwner(),
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	// test will fail if error is nil and account is empty
	require.NoError(t, err)
	require.NotEmpty(t, account) 

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateRandomAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)	
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}


func TestDeleteAccount(t *testing.T){
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccountByID(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.Error(t, err)
	// require.EqualError(t, err, sql.ErrNoRows.Error())
	require.ErrorContains(t, err, "no rows in result set") 
	require.Empty(t, account2)
}


func TestListAccount(t *testing.T){
	// creating 5 random accounts
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)


	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

// A Context carries a deadline, a cancellation signal, and other 
// values across API boundaries.

