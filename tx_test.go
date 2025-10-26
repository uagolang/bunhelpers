package bunhelpers

import (
	"context"
	"errors"
	"testing"
)

func TestTxFromContext(t *testing.T) {
	ctx := context.Background()

	// Test with no transaction
	tx := TxFromContext(ctx)
	if tx != nil {
		t.Error("TxFromContext() should return nil when no transaction in context")
	}

	// Test with transaction
	db := newTestDB()
	defer db.Close()

	bunTx, _ := db.BeginTx(ctx, nil)
	ctx = TxToContext(ctx, &bunTx)

	tx = TxFromContext(ctx)
	if tx == nil {
		t.Error("TxFromContext() should return transaction when one is in context")
	}

	// Test with wrong type in context
	ctx = context.WithValue(context.Background(), TxKey, "not a transaction")
	tx = TxFromContext(ctx)
	if tx != nil {
		t.Error("TxFromContext() should return nil when context value is not *bun.Tx")
	}
}

func TestTxToContext(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	bunTx, _ := db.BeginTx(context.Background(), nil)

	// Test with valid context
	ctx := context.Background()
	newCtx := TxToContext(ctx, &bunTx)

	if newCtx == nil {
		t.Fatal("TxToContext() returned nil context")
	}

	tx := TxFromContext(newCtx)
	if tx == nil {
		t.Error("Transaction should be retrievable from context")
	}

	// Test with nil context
	newCtx = TxToContext(context.Background(), &bunTx)
	if newCtx == nil {
		t.Fatal("TxToContext() should create background context when ctx is nil")
	}

	tx = TxFromContext(newCtx)
	if tx == nil {
		t.Error("Transaction should be retrievable even when nil context was passed")
	}
}

func TestInTx(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	ctx := context.Background()

	t.Run("successful transaction", func(t *testing.T) {
		callCount := 0
		err := InTx(ctx, db, func(ctx context.Context) error {
			callCount++
			tx := TxFromContext(ctx)
			if tx == nil {
				t.Error("Transaction should be available in context")
			}
			return nil
		})

		if err != nil {
			t.Errorf("InTx() returned error: %v", err)
		}

		if callCount != 1 {
			t.Errorf("Function should be called once, got %d", callCount)
		}
	})

	t.Run("transaction with error", func(t *testing.T) {
		testErr := errors.New("test error")
		err := InTx(ctx, db, func(ctx context.Context) error {
			return testErr
		})

		if err == nil {
			t.Error("InTx() should return error when function returns error")
		}

		if !errors.Is(err, testErr) {
			t.Errorf("InTx() returned wrong error: got %v, want %v", err, testErr)
		}
	})

	t.Run("nested transaction", func(t *testing.T) {
		outerCalled := false
		innerCalled := false

		err := InTx(ctx, db, func(outerCtx context.Context) error {
			outerCalled = true
			outerTx := TxFromContext(outerCtx)

			// Nested InTx should use the same transaction
			return InTx(outerCtx, db, func(innerCtx context.Context) error {
				innerCalled = true
				innerTx := TxFromContext(innerCtx)

				if innerTx != outerTx {
					t.Error("Nested transaction should be the same as outer")
				}

				return nil
			})
		})

		if err != nil {
			t.Errorf("Nested InTx() returned error: %v", err)
		}

		if !outerCalled || !innerCalled {
			t.Error("Both outer and inner functions should be called")
		}
	})

	t.Run("panic recovery", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("InTx() should propagate panic")
			}
		}()

		_ = InTx(ctx, db, func(ctx context.Context) error {
			panic("test panic")
		})
	})
}
