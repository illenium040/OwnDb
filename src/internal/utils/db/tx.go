package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func Tx(ctx context.Context, con *pgx.Conn, action func(tx pgx.Tx) error) error {
	tx, err := con.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				err = fmt.Errorf("rollback error: %v, original error: %w", rollbackErr, err)
			}
		}
	}()

	err = action(tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
