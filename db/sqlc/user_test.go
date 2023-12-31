package db

import (
	"context"
	"testing"

	"github.com/imrishuroy/read-cache-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		ID:    util.RandomOwner(),
		Email: util.RandomEmail(),
		Name:  util.RandomOwner(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.ID, user.ID)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Name, user.Name)
	require.NotZero(t, user, user.CreatedAt)

	return user

}
