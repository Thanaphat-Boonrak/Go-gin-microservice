package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransferTxParams struct {
	FromAccountId pgtype.Int8 `json:"from_account_id"`
	ToAccountId   pgtype.Int8 `json:"to_account_id"`
	Amount        int64       `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.Begin(ctx)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx rollback failed: %w", rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccountId.Int64 < arg.ToAccountId.Int64 {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountId.Int64, -arg.Amount, arg.ToAccountId.Int64, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountId.Int64, arg.Amount, arg.FromAccountId.Int64, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalanced(ctx, AddAccountBalancedParams{
		ID:     accountID1,
		Amount: amount1,
	})

	if err != nil {
		return account1, account2, err
	}

	account2, err = q.AddAccountBalanced(ctx, AddAccountBalancedParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
