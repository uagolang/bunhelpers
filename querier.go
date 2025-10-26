package bunhelpers

import (
	"context"

	"github.com/uptrace/bun"
)

type Querier interface {
	NewSelectQuery(ctx context.Context) *bun.SelectQuery
	NewInsertQuery(ctx context.Context) *bun.InsertQuery
	NewUpdateQuery(ctx context.Context) *bun.UpdateQuery
	NewDeleteQuery(ctx context.Context) *bun.DeleteQuery
}

type querier struct {
	db *bun.DB
}

func NewQuerier(c *bun.DB) Querier {
	return &querier{
		db: c,
	}
}

func (r *querier) NewSelectQuery(ctx context.Context) *bun.SelectQuery {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx.NewSelect()
	}
	return r.db.NewSelect()
}

func (r *querier) NewInsertQuery(ctx context.Context) *bun.InsertQuery {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx.NewInsert()
	}
	return r.db.NewInsert()
}

func (r *querier) NewUpdateQuery(ctx context.Context) *bun.UpdateQuery {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx.NewUpdate()
	}
	return r.db.NewUpdate()
}

func (r *querier) NewDeleteQuery(ctx context.Context) *bun.DeleteQuery {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx.NewDelete()
	}
	return r.db.NewDelete()
}
