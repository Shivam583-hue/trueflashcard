package db

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Shivam583-hue/trueflashcard/server/internal/db/dbgen"
)

var ErrMissingDatabaseURL = errors.New("DATABASE_URL is not set")

type Store struct {
	*dbgen.Queries
	pool *pgxpool.Pool
}

func Connect(ctx context.Context) (*Store, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, ErrMissingDatabaseURL
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Store{Queries: dbgen.New(pool), pool: pool}, nil
}

func (s *Store) VerifyConnectivity(ctx context.Context) error {
	_, err := s.Ping(ctx)
	return err
}

func (s *Store) WithTx(ctx context.Context, fn func(q dbgen.Querier) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(s.Queries.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) Close() {
	s.pool.Close()
}
