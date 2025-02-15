package db

import (
	"context"
	"testing"

	"github.com/raihanki/simplebank/util"
	"github.com/stretchr/testify/require"
)

func GenerateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomString(10),
		Password: util.RandomString(8),
		FullName: util.RandomString(6) + " " + util.RandomString(6),
		Email: func() string {
			return util.RandomString(6) + "@gmail.com"
		}(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt)
	require.NotEmpty(t, user.Password)

	return user
}

func TestCreateUser(t *testing.T) {
	GenerateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := GenerateRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.NotZero(t, user2.CreatedAt)
	require.NotZero(t, user2.PasswordChangedAt)
	require.NotEmpty(t, user2.Password)
}
