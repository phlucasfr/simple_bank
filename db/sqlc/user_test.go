package db

import (
	"context"
	"database/sql"
	"picpay_simplificado/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	userParams := CreateUserParams{
		Username:       util.RandomString(5),
		FullName:       util.RandomString(10) + "  " + util.RandomString(10),
		CpfCnpj:        util.RandomCpfCnpj(),
		Email:          util.RandomString(8),
		HashedPassword: util.RandomString(6),
	}

	user, err := testQueries.CreateUser(context.Background(), userParams)
	require.NoError(t, err)

	require.NotEmpty(t, user)
	require.Equal(t, userParams.FullName, user.FullName)
	require.Equal(t, userParams.CpfCnpj, user.CpfCnpj)
	require.Equal(t, userParams.Email, user.Email)
	require.Equal(t, userParams.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.CpfCnpj, user2.CpfCnpj)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	userParams := UpdateUserParams{
		Username:       user1.Username,
		HashedPassword: util.RandomString(12),
		Email:          util.RandomString(6) + "@test.go",
		IsMerchant: sql.NullBool{
			Bool:  false,
			Valid: false,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), userParams)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, userParams.HashedPassword, user2.HashedPassword)
	require.Equal(t, userParams.Email, user2.Email)
	require.Equal(t, userParams.IsMerchant, user2.IsMerchant)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user1.Username)

	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	args := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
