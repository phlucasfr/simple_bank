package db

import (
	"context"
	"database/sql"
	"picpay_simplificado/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomWallet(t *testing.T) Wallet {
	UID := createRandomUser(t).ID

	walletParams := CreateWalletParams{
		UserID:   UID,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	wallet, err := testQueries.CreateWallet(context.Background(), walletParams)

	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	require.Equal(t, walletParams.UserID, wallet.UserID)
	require.Equal(t, walletParams.Balance, wallet.Balance)
	require.Equal(t, walletParams.Currency, wallet.Currency)

	require.NotZero(t, wallet.ID)
	require.NotZero(t, wallet.CreatedAt)

	return wallet
}

func TestCreateWallet(t *testing.T) {
	createRandomWallet(t)
}

func TestGetWallet(t *testing.T) {
	wallet1 := createRandomWallet(t)

	wallet2, err := testQueries.GetWallet(context.Background(), wallet1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.ID, wallet2.ID)
	require.Equal(t, wallet1.UserID, wallet2.UserID)
	require.Equal(t, wallet1.Balance, wallet2.Balance)
	require.Equal(t, wallet1.Currency, wallet2.Currency)

	require.WithinDuration(t, wallet1.CreatedAt.Time, wallet2.CreatedAt.Time, time.Second)
}

func TestUpdateWallet(t *testing.T) {
	wallet1 := createRandomWallet(t)

	walletParams := UpdateWalletParams{
		ID:      wallet1.ID,
		Balance: util.RandomMoney(),
	}

	wallet2, err := testQueries.UpdateWallet(context.Background(), walletParams)

	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.ID, wallet2.ID)
	require.Equal(t, wallet1.UserID, wallet2.UserID)
	require.Equal(t, walletParams.Balance, wallet2.Balance)
	require.Equal(t, wallet1.Currency, wallet2.Currency)

	require.WithinDuration(t, wallet1.CreatedAt.Time, wallet2.CreatedAt.Time, time.Second)
}

func TestDeleteWallet(t *testing.T) {
	wallet1 := createRandomWallet(t)

	err := testQueries.DeleteWallet(context.Background(), wallet1.ID)

	require.NoError(t, err)

	wallet2, err := testQueries.GetWallet(context.Background(), wallet1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, wallet2)
}

func TestListWallets(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomWallet(t)
	}

	args := ListWalletsParams{
		Limit:  5,
		Offset: 5,
	}

	wallets, err := testQueries.ListWallets(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, wallets, 5)

	for _, wallet := range wallets {
		require.NotEmpty(t, wallet)
	}
}
