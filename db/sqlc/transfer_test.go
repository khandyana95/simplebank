package db

import (
	"context"
	"testing"
	"time"

	"github.com/khandyan95/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, toAccount, fromAccount Account) Transfer {

	arg := CreateTransferParams{
		ToAccountID:   toAccount.ID,
		FromAccountID: fromAccount.ID,
		Amount:        util.RandomMoney(),
		Currency:      toAccount.Currency,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.Currency, transfer.Currency)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.Currency, transfer2.Currency)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			createRandomTransfer(t, account1, account2)
		} else {
			createRandomTransfer(t, account2, account1)
		}

	}

	arg := ListTransfersParams{
		ToAccountID:   account1.ID,
		FromAccountID: account1.ID,
		Limit:         10,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 10)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.ToAccountID == account1.ID || transfer.FromAccountID == account1.ID)
	}

	arg.ToAccountID = account1.ID
	arg.FromAccountID = account2.ID

	transfers, err = testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.ToAccountID == account1.ID && transfer.FromAccountID == account2.ID)
	}
}
