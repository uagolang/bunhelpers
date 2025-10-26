package bunhelpers

import (
	"strings"
	"testing"
	"time"

	"github.com/uptrace/bun"
)

func TestWhere_Where(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	type testModel struct {
		bun.BaseModel `bun:"table:test"`
		ID            string `bun:"id,pk"`
		Name          string `bun:"name"`
	}

	t.Run("nil where", func(t *testing.T) {
		var where *Where
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		if result == nil {
			t.Error("Where() should not return nil for nil receiver")
		}
	})

	t.Run("with ID", func(t *testing.T) {
		where := &Where{ID: "test-id"}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "id") {
			t.Error("Where() should add ID condition")
		}
	})

	t.Run("with IDs", func(t *testing.T) {
		where := &Where{IDs: []string{"1", "2", "3"}}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "IN") {
			t.Error("Where() should add IN condition for IDs")
		}
	})

	t.Run("with NotInIDs", func(t *testing.T) {
		where := &Where{NotInIDs: []string{"1", "2"}}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "NOT IN") {
			t.Error("Where() should add NOT IN condition")
		}
	})

	t.Run("with HasFlags", func(t *testing.T) {
		where := &Where{HasFlags: []int{1, 2}}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "flags") {
			t.Error("Where() should add flag conditions")
		}
	})

	t.Run("with HasNotFlags", func(t *testing.T) {
		where := &Where{HasNotFlags: []int{4}}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "flags") {
			t.Error("Where() should add negative flag conditions")
		}
	})

	t.Run("with OnlyDeleted", func(t *testing.T) {
		where := &Where{OnlyDeleted: true}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		// This should call WhereDeleted() on the query
		_ = result // Just ensure it doesn't panic
	})

	t.Run("with WithDeleted", func(t *testing.T) {
		where := &Where{WithDeleted: true}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		// This should call WhereAllWithDeleted() on the query
		_ = result // Just ensure it doesn't panic
	})

	t.Run("with CreatedAfter", func(t *testing.T) {
		now := time.Now().UnixMilli()
		where := &Where{CreatedAfter: &now}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "created_at") {
			t.Error("Where() should add CreatedAfter condition")
		}
	})

	t.Run("with CreatedBefore", func(t *testing.T) {
		now := time.Now().UnixMilli()
		where := &Where{CreatedBefore: &now}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "created_at") {
			t.Error("Where() should add CreatedBefore condition")
		}
	})

	t.Run("with UpdatedAfter", func(t *testing.T) {
		now := time.Now().UnixMilli()
		where := &Where{UpdatedAfter: &now}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "updated_at") {
			t.Error("Where() should add UpdatedAfter condition")
		}
	})

	t.Run("with UpdatedBefore", func(t *testing.T) {
		now := time.Now().UnixMilli()
		where := &Where{UpdatedBefore: &now}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "updated_at") {
			t.Error("Where() should add UpdatedBefore condition")
		}
	})

	t.Run("with custom column names", func(t *testing.T) {
		where := &Where{
			FlagsCol:     "custom_flags",
			CreatedAtCol: "custom_created",
			UpdatedAtCol: "custom_updated",
			HasFlags:     []int{1},
		}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Where(query)

		sql := result.String()
		if !strings.Contains(sql, "flags") {
			t.Error("Where() should use custom column names")
		}
	})
}

func TestWhere_Select(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	type testModel struct {
		bun.BaseModel `bun:"table:test"`
		ID            string `bun:"id,pk"`
		Name          string `bun:"name"`
	}

	t.Run("with SelectColumns", func(t *testing.T) {
		where := Where{SelectColumns: []string{"id", "name"}}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Select(query)

		sql := result.String()
		if !strings.Contains(sql, "id") || !strings.Contains(sql, "name") {
			t.Error("Select() should select specific columns")
		}
	})

	t.Run("with ExcludeColumns", func(t *testing.T) {
		where := Where{ExcludeColumns: []string{"name"}}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Select(query)

		// ExcludeColumn is applied
		_ = result
	})

	t.Run("with Limit", func(t *testing.T) {
		limit := 10
		where := Where{Limit: &limit}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Select(query)

		sql := result.String()
		if !strings.Contains(sql, "LIMIT") {
			t.Error("Select() should add LIMIT")
		}
	})

	t.Run("with Offset", func(t *testing.T) {
		offset := 5
		where := Where{Offset: &offset}
		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Select(query)

		sql := result.String()
		if !strings.Contains(sql, "OFFSET") {
			t.Error("Select() should add OFFSET")
		}
	})

	t.Run("with Order ascending", func(t *testing.T) {
		where := Where{
			SortBy:   1,
			SortDesc: false,
		}
		where.Order = Order{
			1: "name",
		}

		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Select(query)

		sql := result.String()
		if !strings.Contains(sql, "ORDER") || !strings.Contains(sql, "asc") {
			t.Error("Select() should add ORDER BY ascending")
		}
	})

	t.Run("with Order descending", func(t *testing.T) {
		where := Where{
			SortBy:   1,
			SortDesc: true,
		}
		where.Order = Order{
			1: "name",
		}

		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Select(query)

		sql := result.String()
		if !strings.Contains(sql, "ORDER") || !strings.Contains(sql, "desc") {
			t.Error("Select() should add ORDER BY descending")
		}
	})

	t.Run("with invalid SortBy", func(t *testing.T) {
		where := Where{
			SortBy:   999,
			SortDesc: true,
		}
		where.Order = Order{
			1: "name",
		}

		query := db.NewSelect().Model((*testModel)(nil))
		result := where.Select(query)

		// Should not panic with invalid SortBy
		_ = result
	})
}

func TestOrderAsc(t *testing.T) {
	result := OrderAsc("name")
	expected := "name asc"

	if result != expected {
		t.Errorf("OrderAsc() = %q, want %q", result, expected)
	}
}

func TestOrderDesc(t *testing.T) {
	result := OrderDesc("name")
	expected := "name desc"

	if result != expected {
		t.Errorf("OrderDesc() = %q, want %q", result, expected)
	}
}

func TestWhereConstants(t *testing.T) {
	if DefaultIDCol != "id" {
		t.Errorf("DefaultIDCol = %q, want %q", DefaultIDCol, "id")
	}
	if DefaultFlagsCol != "flags" {
		t.Errorf("DefaultFlagsCol = %q, want %q", DefaultFlagsCol, "flags")
	}
	if DefaultCreatedAtCol != "created_at" {
		t.Errorf("DefaultCreatedAtCol = %q, want %q", DefaultCreatedAtCol, "created_at")
	}
	if DefaultUpdatedAtCol != "updated_at" {
		t.Errorf("DefaultUpdatedAtCol = %q, want %q", DefaultUpdatedAtCol, "updated_at")
	}
}

func TestOrder(t *testing.T) {
	order := Order{
		1: "created_at",
		2: "updated_at",
		3: "name",
	}

	if order[1] != "created_at" {
		t.Error("Order map should store values correctly")
	}
	if order[2] != "updated_at" {
		t.Error("Order map should store values correctly")
	}
	if order[3] != "name" {
		t.Error("Order map should store values correctly")
	}
}
