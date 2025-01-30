package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {

	password := randomString(10)
	hashedPassword1, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = VerifyPassword(password, hashedPassword1)

	require.NoError(t, err)

	err = VerifyPassword(randomString(10), hashedPassword1)
	require.Error(t, err, bcrypt.ErrMismatchedHashAndPassword)

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
