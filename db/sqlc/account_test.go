package db

import (
	"context"
	"testing"
	"time"

	"github.com/raihanki/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := GenerateRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomInt(0, 1000),
		Currency: "USD",
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CraetedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	myAccount := createRandomAccount(t)
	account, err := testQueries.GetAccount(context.Background(), myAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, myAccount.ID, account.ID)
	require.Equal(t, myAccount.Owner, account.Owner)
	require.Equal(t, myAccount.Balance, account.Balance)
	require.Equal(t, myAccount.Currency, account.Currency)
	require.WithinDuration(t, myAccount.CraetedAt.Time, account.CraetedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	myAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      myAccount.ID,
		Balance: util.RandomInt(0, 1000),
	}
	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, myAccount.ID, account.ID)
	require.Equal(t, myAccount.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, myAccount.Currency, account.Currency)
	require.WithinDuration(t, myAccount.CraetedAt.Time, account.CraetedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	myAccount := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), myAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), myAccount.ID)
	require.Error(t, err)
	// require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestGetAllAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := GetAllAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.GetAllAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
