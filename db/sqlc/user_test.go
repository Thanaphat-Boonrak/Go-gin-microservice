package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/thanaphat2005/simple-bank/utils"
)

func createRandomUser(t *testing.T) User {

	hashPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)
	arg := CreateUsersParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandonEmail(),
	}

	user, err := testQueries.CreateUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUsers(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)

}
