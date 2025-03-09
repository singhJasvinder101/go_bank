package db

import (
	"context"
	"testing"

	"github.com/singhJasvinder101/go_bank/utils"
	"github.com/stretchr/testify/require"
)

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


