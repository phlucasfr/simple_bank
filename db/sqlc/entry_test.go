package db

import (
	"context"
	"picpay_simplificado/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, wallet Wallet) Entry {
	entryParams := CreateEntryParams{
		WalletID: wallet.ID,
		Amount:   util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), entryParams)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entryParams.WalletID, entry.WalletID)
	require.Equal(t, entryParams.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	wallet := createRandomWallet(t)
	createRandomEntry(t, wallet)
}

func TestGetEntry(t *testing.T) {
	wallet := createRandomWallet(t)
	entry1 := createRandomEntry(t, wallet)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.WalletID, entry2.WalletID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	wallet := createRandomWallet(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, wallet)
	}

	arg := ListEntriesParams{
		WalletID: wallet.ID,
		Limit:    5,
		Offset:   5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.WalletID, entry.WalletID)
	}
}
