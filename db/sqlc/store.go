package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return Store{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// execTx executes a function within a database transaction.
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// Create the transfer record.
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		// Create the from account entry.
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create the to account entry.
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update account balances.
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
