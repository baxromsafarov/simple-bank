package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple-bank/util"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomString(6) + "@email.com",
	}

	// Используем NewStore для получения Queries
	store := NewStore(testDB)
	user, err := store.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  1000,
		Currency: "USD",
	}

	// Используем NewStore для получения Queries
	store := NewStore(testDB)
	account, err := store.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	return account
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  1000,
		Currency: "USD",
	}
	store := NewStore(testDB)
	account, err := store.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
