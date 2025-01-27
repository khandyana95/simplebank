package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	acc1Balance := account1.Balance
	acc2Balance := account2.Balance

	arg := CreateTransactionParams{
		FromAccountId: account1.ID,
		ToAccountId:   account2.ID,
		Amount:        10,
		Currency:      account1.Currency,
	}

	errChan := make(chan error)
	txnResultChan := make(chan TransferTxnResult)

	for i := 0; i < 5; i++ {
		go func() {
			txnResult, err := testStore.CreateTransaction(context.Background(), arg)

			errChan <- err
			txnResultChan <- txnResult
		}()
	}

	for i := 0; i < 5; i++ {

		err := <-errChan
		require.NoError(t, err)

		txnResult := <-txnResultChan
		require.NotEmpty(t, txnResult)

		transfer := txnResult.TransferTxn
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.CreatedAt)
		require.Equal(t, arg.Amount, transfer.Amount)
		require.Equal(t, arg.Currency, transfer.Currency)

		_, err1 := testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err1)

		entry := txnResult.FromAccountEntry
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.Equal(t, account1.ID, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
		require.Equal(t, -arg.Amount, entry.Amount)
		require.Equal(t, arg.Currency, entry.Currency)

		_, err1 = testQueries.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err1)

		entry = txnResult.ToAccountEntry
		require.NotEmpty(t, entry)
		require.NotZero(t, entry.ID)
		require.Equal(t, account2.ID, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
		require.Equal(t, arg.Amount, entry.Amount)
		require.Equal(t, arg.Currency, entry.Currency)

		_, err1 = testQueries.GetEntry(context.Background(), entry.ID)
		require.NoError(t, err1)

		fromAccount := txnResult.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		require.Equal(t, account1.Name, fromAccount.Name)
		require.Equal(t, account1.Owner, fromAccount.Owner)
		require.Equal(t, account1.Currency, fromAccount.Currency)
		amount := acc1Balance - fromAccount.Balance
		require.Equal(t, amount, float64(i+1)*arg.Amount)
		fmt.Println("from account balance", fromAccount.Balance)

		toAccount := txnResult.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		require.Equal(t, account2.Name, toAccount.Name)
		require.Equal(t, account2.Owner, toAccount.Owner)
		require.Equal(t, account2.Currency, toAccount.Currency)
		amount = toAccount.Balance - acc2Balance
		require.Equal(t, amount, float64(i+1)*arg.Amount)
	}
}
