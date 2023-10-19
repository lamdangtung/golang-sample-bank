package db

import (
	"context"
	"testing"

	"github.com/lamdangtung/golang-sample-bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func CreateRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashedPassword(util.RandomPassword())
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		FullName:       util.RandomFullname(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.Equal(t, user1.FullName, user2.FullName)

}
