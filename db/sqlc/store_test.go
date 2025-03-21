package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)
func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	// queue channels for errors and results
	errs := make(chan error)
	results := make(chan TransferTxResult)


	fmt.Println("Before: ", account1.Balance, account2.Balance)
	for i := 0; i < n; i++ {
		go func() {
			txName := fmt.Sprintf("tx %d", i+1)

			// ctxWithVal := context.WithValue(context.Background(), txKey, txName)
			ctx := context.WithValue(context.Background(), txKey, txName)
			// ctx, cancel := context.WithTimeout(ctxWithVal, time.Second)
			// defer cancel()


			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, amount, transfer.CreatedAt)

		// 1 check transfer
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// 2 check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// 3 check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// 4 check account balances
		fmt.Println(">>Tx ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)	
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)

	// 200 - n (amount) after n concurrent transactions
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	// queue channels for errors and results
	errs := make(chan error)

	fmt.Println("Before: ", account1.Balance, account2.Balance)
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId, toAccountId = toAccountId, fromAccountId
		}

		go func() {
			txName := fmt.Sprintf("tx %d", i+1)

			// ctxWithVal := context.WithValue(context.Background(), txKey, txName)
			ctx := context.WithValue(context.Background(), txKey, txName)
			// ctx, cancel := context.WithTimeout(ctxWithVal, time.Second)
			// defer cancel()

			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)	
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)

	// Now balance must same as before
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
