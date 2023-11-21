package api

import (
	db "picpay_simplificado/db/sqlc"
	"picpay_simplificado/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string, err error) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomString(5),
		FullName:       util.RandomString(10) + " " + util.RandomString(10),
		CpfCnpj:        util.RandomCpfCnpj(),
		Email:          user.FullName + "@test.go",
		HashedPassword: hashedPassword,
	}
	return
}

func TestCreateUser(t *testing.T) {
	user, password, err := randomUser(t)
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Username:       user.Username,
		FullName:       user.FullName,
		CpfCnpj:        user.CpfCnpj,
		Email:          user.Email,
		HashedPassword: password,
	}
	require.NotEmpty(t, arg)

	//TODO: implementar testes
}
