package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrasferTx(t *testing.T) {
	store := NewStore(testDB)

	wallet1 := createRandomWallet(t)
	wallet2 := createRandomWallet(t)
	fmt.Println("DEBUG>> Before wallet1 balance: ", wallet1.Balance, "wallet2 balance: ", wallet2.Balance)

	n := 5

	amountFloat := 10.00
	amount := int64(amountFloat * 100)

	errs := make(chan error)
	results := make(chan TrasferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TrasferTxParms{
				FromWalletID: wallet1.ID,
				ToWalletID:   wallet2.ID,
				Amount:       amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, wallet1.ID, transfer.FromWalletID)
		require.Equal(t, wallet2.ID, transfer.ToWalletID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)

		require.Equal(t, wallet1.ID, fromEntry.WalletID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)

		require.Equal(t, wallet2.ID, toEntry.WalletID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check wallets
		fromWallet := result.FromWallet
		require.NotEmpty(t, fromWallet)
		require.Equal(t, wallet1.ID, fromWallet.ID)

		toWallet := result.ToWallet
		require.NotEmpty(t, toWallet)
		require.Equal(t, wallet2.ID, toWallet.ID)

		//check wallets balance
		fmt.Println("DEBUG>> During transaction wallet1 balance: ", fromWallet.Balance, "wallet2 balance: ", toWallet.Balance)
		diff1 := wallet1.Balance - fromWallet.Balance
		diff2 := toWallet.Balance - wallet2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)  //number of transactions
		require.NotContains(t, existed, k) //existed is out of loop
		existed[k] = true
	}

	//check the final updates balance
	updatedWallet1, err := store.GetWallet(context.Background(), wallet1.ID)
	require.NoError(t, err)

	updatedWallet2, err := store.GetWallet(context.Background(), wallet2.ID)
	require.NoError(t, err)

	fmt.Println("DEBUG>> After wallet1 balance: ", updatedWallet1.Balance, "wallet2 balance: ", updatedWallet2.Balance)
	require.Equal(t, wallet1.Balance-int64(n)*amount, updatedWallet1.Balance)
	require.Equal(t, wallet2.Balance+int64(n)*amount, updatedWallet2.Balance)
}

func TestTrasferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	wallet1 := createRandomWallet(t)
	wallet2 := createRandomWallet(t)
	fmt.Println("DEBUG>> Before wallet1 balance: ", wallet1.Balance, "wallet2 balance: ", wallet2.Balance)

	n := 10

	amountFloat := 10.00
	amount := int64(amountFloat * 100)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromWalletID := wallet1.ID
		toWalletID := wallet2.ID

		if i%2 == 1 {
			fromWalletID = wallet2.ID
			toWalletID = wallet1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TrasferTxParms{
				FromWalletID: fromWalletID,
				ToWalletID:   toWalletID,
				Amount:       amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	//check the final updates balance
	updatedWallet1, err := store.GetWallet(context.Background(), wallet1.ID)
	require.NoError(t, err)

	updatedWallet2, err := store.GetWallet(context.Background(), wallet2.ID)
	require.NoError(t, err)

	fmt.Println("DEBUG>> After wallet1 balance: ", updatedWallet1.Balance, "wallet2 balance: ", updatedWallet2.Balance)
	require.Equal(t, wallet1.Balance, updatedWallet1.Balance)
	require.Equal(t, wallet2.Balance, updatedWallet2.Balance)
}
