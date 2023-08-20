package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	transaction, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	querie := New(transaction)
	err = fn(querie)

	if err != nil {
		if rbError := transaction.Rollback(); rbError != nil {
			return fmt.Errorf("transaction error: %v , rollback error: %v", err, rbError)
		}

		return err
	}

	return transaction.Commit()
}

type TrasferTxParms struct {
	FromWalletID int64 `json:"from_wallet_id"`
	ToWalletID   int64 `json:"to_wallet_id"`
	Amount       int64 `json:"amount"`
}

type TrasferTxResult struct {
	Transfer   Transfer `json:"transfer"`
	FromWallet Wallet   `json:"from_wallet"`
	ToWallet   Wallet   `json:"to_wallet"`
	FromEntry  Entry    `json:"from_entry"`
	ToEntry    Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TrasferTxParms) (TrasferTxResult, error) {
	var result TrasferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromWalletID: arg.FromWalletID,
			ToWalletID:   arg.ToWalletID,
			Amount:       arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			WalletID: arg.FromWalletID,
			Amount:   -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			WalletID: arg.ToWalletID,
			Amount:   arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromWalletID < arg.ToWalletID {
			result.FromWallet, result.ToWallet, err = addMoney(ctx, q, arg.FromWalletID, -arg.Amount, arg.ToWalletID, arg.Amount)
		} else {
			result.ToWallet, result.FromWallet, err = addMoney(ctx, q, arg.ToWalletID, arg.Amount, arg.FromWalletID, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	walletID1 int64,
	amount1 int64,
	walletID2 int64,
	amount2 int64,
) (wallet1 Wallet, wallet2 Wallet, err error) {
	wallet1, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
		Amount: amount1,
		ID:     walletID1,
	})
	if err != nil {
		return
	}

	wallet2, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
		Amount: amount2,
		ID:     walletID2,
	})

	return
}
