package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashedPass(t *testing.T) {
	password := RandomPassword()

	hashedPassword, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomPassword()
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
