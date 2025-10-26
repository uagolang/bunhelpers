package bunhelpers

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type txKey int

const TxKey txKey = 1

func TxFromContext(ctx context.Context) *bun.Tx {
	tx, ok := ctx.Value(TxKey).(*bun.Tx)
	if !ok {
		return nil
	}
	return tx
}

func TxToContext(ctx context.Context, tx *bun.Tx) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, TxKey, tx)
}

func InTx(ctx context.Context, client *bun.DB, fn func(ctx context.Context) error) error {
	var err error
	var rootTx bool

	tx := TxFromContext(ctx)
	if tx == nil {
		rootTx = true

		_tx, err := client.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		tx = &_tx
	}

	ctxWithTx := TxToContext(ctx, tx)

	if !rootTx {
		return fn(ctxWithTx)
	}

	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	err = fn(ctxWithTx)

	if err == nil {
		err := tx.Commit()
		if err != nil {
			return err
		}

		return nil
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		err = fmt.Errorf("%w: transaction rollback error: %v", err, rollbackErr)
	}
	return err
}
