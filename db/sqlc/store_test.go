package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDb)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	parallelRuns := 5
	amount := int64(10)

	//fmt.Println("Accounts Before:")
	//fmt.Printf("FromAccount: %+v\n", fromAccount)
	//fmt.Printf("ToAccount: %+v\n\n", toAccount)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < parallelRuns; i++ {
		txName := fmt.Sprintf("tx: %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < parallelRuns; i++ {
		err := <-errs

		require.NoError(t, err)

		result := <-results

		//fmt.Println("Transfer")
		//fmt.Printf("transfer: %+v\n\n", result.Transfer)
		//fmt.Printf("FromEntry: %+v\n", result.FromEntry)
		//fmt.Printf("ToEntry: %+v\n\n", result.ToEntry)
		//fmt.Printf("FromAccount: %+v\n", result.FromAccount)
		//fmt.Printf("ToAccount: %+v\n\n", result.ToAccount)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, transfer.ToAccountID, toAccount.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		toEntry := result.ToEntry
		fromEntry := result.FromEntry

		//require.NotEmpty(t, toEntry)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, fromAccount.ID)
		require.Equal(t, toEntry.AccountID, toAccount.ID)
		require.Equal(t, toEntry.Amount, -fromEntry.Amount)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// test accounts
		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, result.ToAccount.ID, toAccount.ID)

		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, result.FromAccount.ID, fromAccount.ID)

		fromDiff := fromAccount.Balance - result.FromAccount.Balance
		toDiff := result.ToAccount.Balance - toAccount.Balance
		require.Equal(t, fromDiff, toDiff)

		require.True(t, fromDiff > 0)
		require.True(t, fromDiff%amount == 0) // is multiple of amount

		k := int(fromDiff / amount)
		require.True(t, k >= 1 && k <= parallelRuns)

		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final balance
	totalAmount := int64(parallelRuns) * amount
	fromAccountAfter, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	require.Equal(t, fromAccount.Balance-totalAmount, fromAccountAfter.Balance)

	toAccountAfter, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)
	require.Equal(t, toAccount.Balance+totalAmount, toAccountAfter.Balance)

}

func TestStore_TransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	parallelRuns := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < parallelRuns; i++ {
		txName := fmt.Sprintf("tx: %d", i+1)
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 0 {
			fromAccountID, toAccountID = toAccountID, fromAccountID
		}

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < parallelRuns; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check final balance
	account1After, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1After.Balance, account1.Balance)

	account2After, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2After.Balance, account2.Balance)
}
