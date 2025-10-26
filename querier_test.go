package bunhelpers

import (
	"context"
	"testing"
)

func TestNewQuerier(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	querier := NewQuerier(db)
	if querier == nil {
		t.Fatal("NewQuerier() returned nil")
	}

	// Check interface satisfaction
	var _ Querier = querier
}

func TestQuerier_NewSelectQuery(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	querier := NewQuerier(db)
	ctx := context.Background()

	t.Run("without transaction", func(t *testing.T) {
		query := querier.NewSelectQuery(ctx)
		if query == nil {
			t.Fatal("NewSelectQuery() returned nil")
		}
	})

	t.Run("with transaction", func(t *testing.T) {
		bunTx, _ := db.BeginTx(ctx, nil)
		txCtx := TxToContext(ctx, &bunTx)

		query := querier.NewSelectQuery(txCtx)
		if query == nil {
			t.Fatal("NewSelectQuery() returned nil with transaction")
		}
	})
}

func TestQuerier_NewInsertQuery(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	querier := NewQuerier(db)
	ctx := context.Background()

	t.Run("without transaction", func(t *testing.T) {
		query := querier.NewInsertQuery(ctx)
		if query == nil {
			t.Fatal("NewInsertQuery() returned nil")
		}
	})

	t.Run("with transaction", func(t *testing.T) {
		bunTx, _ := db.BeginTx(ctx, nil)
		txCtx := TxToContext(ctx, &bunTx)

		query := querier.NewInsertQuery(txCtx)
		if query == nil {
			t.Fatal("NewInsertQuery() returned nil with transaction")
		}
	})
}

func TestQuerier_NewUpdateQuery(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	querier := NewQuerier(db)
	ctx := context.Background()

	t.Run("without transaction", func(t *testing.T) {
		query := querier.NewUpdateQuery(ctx)
		if query == nil {
			t.Fatal("NewUpdateQuery() returned nil")
		}
	})

	t.Run("with transaction", func(t *testing.T) {
		bunTx, _ := db.BeginTx(ctx, nil)
		txCtx := TxToContext(ctx, &bunTx)

		query := querier.NewUpdateQuery(txCtx)
		if query == nil {
			t.Fatal("NewUpdateQuery() returned nil with transaction")
		}
	})
}

func TestQuerier_NewDeleteQuery(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	querier := NewQuerier(db)
	ctx := context.Background()

	t.Run("without transaction", func(t *testing.T) {
		query := querier.NewDeleteQuery(ctx)
		if query == nil {
			t.Fatal("NewDeleteQuery() returned nil")
		}
	})

	t.Run("with transaction", func(t *testing.T) {
		bunTx, _ := db.BeginTx(ctx, nil)
		txCtx := TxToContext(ctx, &bunTx)

		query := querier.NewDeleteQuery(txCtx)
		if query == nil {
			t.Fatal("NewDeleteQuery() returned nil with transaction")
		}
	})
}
