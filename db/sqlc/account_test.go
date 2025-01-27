package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/khandyan95/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Name:     util.RandomName(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Name, acc.Name)
	require.Equal(t, arg.Currency, acc.Currency)
	require.Equal(t, arg.Balance, acc.Balance)

	require.NotZero(t, acc.CreatedAt)
	require.NotZero(t, acc.ID)

	return acc
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Name, acc2.Name)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Name, acc2.Name)
	// require.Equal(t, arg.Balance, acc2.Balance)
	require.Equal(t, arg.Balance+acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {

	acc1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Empty(t, acc2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 3,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)
	// require.Equal(t, 5, len(accounts))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
