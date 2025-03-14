package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries		 // methods provided by sqlc generated Queries struct
	db *pgxpool.Pool // connection pool for PSQL to begin db.BeginTx
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,	
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	// new Queries object tied for each transaction (tx) to execute 
	// db operations within that specific transaction’s scope, 
	// ensuring atomicity and isolation
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
    var result TransferTxResult

    err := store.execTx(ctx, func(q *Queries) error {
        var err error
		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
        result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
        if err != nil {
            return err
        }

		fmt.Println(txName, "create entry 1")
        result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.FromAccountID,
            Amount:    -arg.Amount,
        })
        if err != nil {
            return err
        }

		fmt.Println(txName, "create entry 2")
        result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.ToAccountID,
            Amount:    arg.Amount,
        })
        if err != nil {
            return err
        }

		fmt.Println(txName, "get account 1: for update")
        // Use FOR UPDATE to lock the rows
        account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
        if err != nil {
            return err
        }

		fmt.Println(txName, "get account 2: for update")
        account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
        if err != nil {
            return err
        }

		fmt.Println(txName, "update account 1")
        result.FromAccount, err = q.UpdateAccountByID(ctx, UpdateAccountByIDParams{
            ID:      arg.FromAccountID,
            Balance: account1.Balance - arg.Amount,
        })
        if err != nil {
            return err
        }

		fmt.Println(txName, "update account 2")
        result.ToAccount, err = q.UpdateAccountByID(ctx, UpdateAccountByIDParams{
            ID:      arg.ToAccountID,
            Balance: account2.Balance + arg.Amount,
        })
        if err != nil {
            return err
        }

        return nil
    })

    return result, err
}
