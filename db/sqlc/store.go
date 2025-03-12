package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// provide all functions to execute db queries and transactions
type Store struct{
	*Queries
	db *pgxpool.Pool
}

// create new store
func NewStore(db *pgxpool.Pool) *Store{
	return &Store {
		db: db,
		Queries: New(db),  
		// return Query object with db operation methods
	}
}

// execTx executes a function within a database transaction
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil{
		return err
	}
	
	// return new object queries for tx
	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
