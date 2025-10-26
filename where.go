package bunhelpers

import (
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

const (
	DefaultIDCol        = "id"
	DefaultFlagsCol     = "flags"
	DefaultCreatedAtCol = "created_at"
	DefaultUpdatedAtCol = "updated_at"
)

type Where struct {
	ID       string   `json:"id,omitempty" form:"id"`
	IDs      []string `json:"ids,omitempty" form:"ids"`
	NotInIDs []string `json:"not_in_ids,omitempty" form:"not_in_ids"`

	HasFlags    []int `json:"has_flags,omitempty" form:"has_flags"`
	HasNotFlags []int `json:"has_not_flags,omitempty" form:"has_not_flags"`

	OnlyDeleted bool `json:"only_deleted,omitempty" form:"only_deleted"`
	WithDeleted bool `json:"with_deleted,omitempty" form:"with_deleted"`

	Limit  *int `json:"limit,omitempty" form:"limit"`
	Offset *int `json:"offset,omitempty" form:"offset"`

	FlagsCol     string `json:"flags_col,omitempty" form:"flags_col"`
	CreatedAtCol string `json:"created_at_col,omitempty" form:"created_at_col"`
	UpdatedAtCol string `json:"updated_at_col,omitempty" form:"updated_at_col"`

	CreatedAfter  *int64 `json:"created_after" form:"created_after"`
	CreatedBefore *int64 `json:"created_before" form:"created_before"`

	UpdatedAfter  *int64 `json:"updated_after" form:"updated_after"`
	UpdatedBefore *int64 `json:"updated_before" form:"updated_before"`

	SelectColumns  []string `json:"select_columns,omitempty" form:"select_columns"`
	ExcludeColumns []string `json:"exclude_columns,omitempty" form:"exclude_columns"`

	SortBy   int  `json:"sort_by,omitempty" form:"sort_by"`
	SortDesc bool `json:"sort_desc,omitempty" form:"sort_desc"`

	Order Order `json:"-"`
}

func (w *Where) Where(q *bun.SelectQuery) *bun.SelectQuery {
	if w == nil {
		return q
	}

	if w.FlagsCol == "" {
		w.FlagsCol = DefaultFlagsCol
	}
	if w.CreatedAtCol == "" {
		w.CreatedAtCol = DefaultCreatedAtCol
	}
	if w.UpdatedAtCol == "" {
		w.UpdatedAtCol = DefaultUpdatedAtCol
	}

	if w.ID != "" {
		q = q.Where("?TableAlias.? = ?", bun.Ident(DefaultIDCol), w.ID)
	}
	if len(w.IDs) > 0 {
		q = q.Where("?TableAlias.? IN (?)", bun.Ident(DefaultIDCol), bun.In(w.IDs))
	}
	if len(w.NotInIDs) > 0 {
		q = q.Where("?TableAlias.? NOT IN (?)", bun.Ident(DefaultIDCol), bun.In(w.NotInIDs))
	}

	for _, flag := range w.HasFlags {
		q = q.Where("?TableAlias.? & ? = ?", bun.Ident(DefaultFlagsCol), flag, flag)
	}
	for _, flag := range w.HasNotFlags {
		q = q.Where("?TableAlias.? & ? = 0", bun.Ident(DefaultFlagsCol), flag)
	}

	if w.OnlyDeleted {
		q.WhereDeleted()
	} else if w.WithDeleted {
		q.WhereAllWithDeleted()
	}

	if w.CreatedAfter != nil {
		q.Where("?TableAlias.? >= ?", bun.Ident(DefaultCreatedAtCol), time.UnixMilli(*w.CreatedAfter))
	}
	if w.CreatedBefore != nil {
		q.Where("?TableAlias.? <= ?", bun.Ident(DefaultCreatedAtCol), time.UnixMilli(*w.CreatedBefore))
	}

	if w.UpdatedAfter != nil {
		q.Where("?TableAlias.? >= ?", bun.Ident(DefaultUpdatedAtCol), time.UnixMilli(*w.UpdatedAfter))
	}
	if w.UpdatedBefore != nil {
		q.Where("?TableAlias.? <= ?", bun.Ident(DefaultUpdatedAtCol), time.UnixMilli(*w.UpdatedBefore))
	}

	return q
}

func (w *Where) Select(q *bun.SelectQuery) *bun.SelectQuery {
	if w == nil {
		return q
	}

	if len(w.SelectColumns) > 0 {
		q = q.Column(w.SelectColumns...)
	}
	if len(w.ExcludeColumns) > 0 {
		q = q.ExcludeColumn(w.ExcludeColumns...)
	}

	if w.Limit != nil {
		q = q.Limit(*w.Limit)
	}
	if w.Offset != nil {
		q = q.Offset(*w.Offset)
	}

	if col, ok := w.Order[w.SortBy]; ok {
		if w.SortDesc {
			q = q.Order(OrderDesc(col))
		} else {
			q = q.Order(OrderAsc(col))
		}
	}

	return q
}

type Order map[int]string

func OrderAsc(col string) string {
	return fmt.Sprintf("%s asc", col)
}

func OrderDesc(col string) string {
	return fmt.Sprintf("%s desc", col)
}
