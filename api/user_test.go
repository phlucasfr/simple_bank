package api

import (
	db "picpay_simplificado/db/sqlc"
	"picpay_simplificado/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		FullName: util.RandomString(10),
		CpfCnpj:  util.RandomCpfCnpj(),
		Email:    user.FullName + "@test.go",
		Password: hashedPassword,
	}
	return
}
