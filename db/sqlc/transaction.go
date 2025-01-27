package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type Store struct {
	*Queries
	DB *sql.DB
}

func NewStore(dbConn *sql.DB) *Store {
	return &Store{
		Queries: New(dbConn),
		DB:      dbConn,
	}
}

var m sync.Mutex

func (store *Store) createDBTxn(ctx context.Context, f func(*Queries) error) error {

	txn, err := store.DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	queryTxn := store.WithTx(txn)

	err = f(queryTxn)

	if err != nil {
		if rbErr := txn.Rollback(); rbErr != nil {
			return fmt.Errorf("txn error %v, rollback error %v", err.Error(), rbErr.Error())
		}

		return err
	}

	return txn.Commit()
}

type CreateTransactionParams struct {
	FromAccountId int64   `json:"from_account_id"`
	ToAccountId   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
}

type TransferTxnResult struct {
	FromAccount      Account  `json:"from_account"`
	ToAccount        Account  `json:"to_account"`
	FromAccountEntry Entry    `json:"from_account_entry"`
	ToAccountEntry   Entry    `json:"to_account_entry"`
	TransferTxn      Transfer `json:"transfer_txn"`
}

func (store *Store) CreateTransaction(ctx context.Context, txnParams CreateTransactionParams) (TransferTxnResult, error) {

	var transferTxnResult TransferTxnResult

	m.Lock()

	err := store.createDBTxn(ctx, func(q *Queries) error {

		fromAccountEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: txnParams.FromAccountId,
			Amount:    -txnParams.Amount,
			Currency:  txnParams.Currency,
		})

		if err != nil {
			return err
		}

		toAccountEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: txnParams.ToAccountId,
			Amount:    txnParams.Amount,
			Currency:  txnParams.Currency,
		})

		if err != nil {
			return err
		}

		fromAccount, err := q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      txnParams.FromAccountId,
			Balance: -txnParams.Amount,
		})

		// fromAccount, err := q.GetAccount(ctx, txnParams.FromAccountId)
		// if err != nil {
		// 	return err
		// }

		// fromAccUpdate := UpdateAccountParams{
		// 	ID:      txnParams.FromAccountId,
		// 	Balance: fromAccount.Balance - txnParams.Amount,
		// }

		// fromAccount, err = q.UpdateAccount(ctx, fromAccUpdate)
		if err != nil {
			return err
		}

		toAccount, err := q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      txnParams.ToAccountId,
			Balance: txnParams.Amount,
		})

		// toAccount, err := q.GetAccount(ctx, txnParams.ToAccountId)
		// if err != nil {
		// 	return err
		// }

		// toAccUpdate := UpdateAccountParams{
		// 	ID:      txnParams.ToAccountId,
		// 	Balance: toAccount.Balance + txnParams.Amount,
		// }

		// toAccount, err = q.UpdateAccount(ctx, toAccUpdate)
		if err != nil {
			return err
		}

		transfer, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: txnParams.FromAccountId,
			ToAccountID:   txnParams.ToAccountId,
			Amount:        txnParams.Amount,
			Currency:      txnParams.Currency,
		})

		if err != nil {
			return err
		}

		transferTxnResult.FromAccount = fromAccount
		transferTxnResult.ToAccount = toAccount
		transferTxnResult.FromAccountEntry = fromAccountEntry
		transferTxnResult.ToAccountEntry = toAccountEntry
		transferTxnResult.TransferTxn = transfer

		return nil
	})

	m.Unlock()

	return transferTxnResult, err
}
