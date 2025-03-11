package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// provide all functions to execute db queries and transactions
type Store struct{
	*Queries
	db *pgxpool.Pool
}

// create new store
func NewStore(db *pgxpool.Pool) *Store{
	return &Store{
		db: db,
		Queries: New(db),  
		// return Query object with db operation methods
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil{
		return err
	}
	
	q := New(tx)
	err = fn(q)
	if err != nil{
		_ = tx.Rollback(ctx)
	}

}