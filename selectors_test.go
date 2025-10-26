package bunhelpers

import (
	"strings"
	"testing"
	"time"

	"github.com/uptrace/bun"
)

type testModel struct {
	bun.BaseModel `bun:"table:test"`
	ID            string `bun:"id,pk"`
	Name          string `bun:"name"`
}

func TestApply(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))

	selector := Apply(
		WhereEqual("name", "test"),
		WhereNotNull("id"),
	)

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") {
		t.Error("Apply() should apply selectors")
	}
}

func TestApplyIf(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	t.Run("condition true", func(t *testing.T) {
		query := db.NewSelect().Model((*testModel)(nil))
		selector := ApplyIf(true, WhereEqual("name", "test"))

		if selector == nil {
			t.Fatal("ApplyIf() should return selector when condition is true")
		}

		result := selector(query)
		sql := result.String()

		if !strings.Contains(sql, "name") {
			t.Error("ApplyIf() should apply selector when condition is true")
		}
	})

	t.Run("condition false", func(t *testing.T) {
		selector := ApplyIf(false, WhereEqual("name", "test"))

		if selector != nil {
			t.Error("ApplyIf() should return nil when condition is false")
		}
	})
}

func TestOrGroup(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))

	// Add a base condition first so OR can be added
	query = WhereEqual("id", "1")(query)

	selector := OrGroup(
		WhereEqual("name", "test1"),
		WhereEqual("name", "test2"),
	)

	result := selector(query)
	sql := result.String()

	// Should have the OR group condition
	if !strings.Contains(sql, "name") {
		t.Error("OrGroup() should add name conditions")
	}
}

func TestAndGroup(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := AndGroup(
		WhereEqual("name", "test1"),
		WhereNotNull("id"),
	)

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") {
		t.Error("AndGroup() should create AND group")
	}
}

func TestOr(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := Or(
		WhereEqual("name", "test1"),
		WhereEqual("name", "test2"),
		WhereEqual("name", "test3"),
	)

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") {
		t.Error("Or() should apply multiple OR conditions")
	}
}

func TestWhereEqual(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereEqual("name", "test")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") || !strings.Contains(sql, "=") {
		t.Error("WhereEqual() should add equality condition")
	}
}

func TestWhereNotEqual(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereNotEqual("name", "test")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") || !strings.Contains(sql, "!=") {
		t.Error("WhereNotEqual() should add not equal condition")
	}
}

func TestWhereNull(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereNull("name")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") || !strings.Contains(sql, "is null") {
		t.Error("WhereNull() should add IS NULL condition")
	}
}

func TestWhereNotNull(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereNotNull("name")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") || !strings.Contains(sql, "is not null") {
		t.Error("WhereNotNull() should add IS NOT NULL condition")
	}
}

func TestWhereIn(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereIn("id", []string{"1", "2", "3"})

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "id") || !strings.Contains(sql, "IN") {
		t.Error("WhereIn() should add IN condition")
	}
}

func TestWhereNotIn(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereNotIn("id", []string{"1", "2", "3"})

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "id") || !strings.Contains(sql, "NOT IN") {
		t.Error("WhereNotIn() should add NOT IN condition")
	}
}

func TestWhereContains(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereContains("name", "test")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") || !strings.Contains(sql, "ILIKE") {
		t.Error("WhereContains() should add ILIKE condition")
	}
}

func TestWhereBegins(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereBegins("name", "test")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") || !strings.Contains(sql, "ILIKE") {
		t.Error("WhereBegins() should add ILIKE condition")
	}
}

func TestWhereEnds(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereEnds("name", "test")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "name") || !strings.Contains(sql, "ILIKE") {
		t.Error("WhereEnds() should add ILIKE condition")
	}
}

func TestWhereBefore(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	now := time.Now()
	selector := WhereBefore("created_at", now)

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "created_at") || !strings.Contains(sql, "<=") {
		t.Error("WhereBefore() should add <= condition")
	}
}

func TestWhereAfter(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	now := time.Now()
	selector := WhereAfter("created_at", now)

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "created_at") || !strings.Contains(sql, ">=") {
		t.Error("WhereAfter() should add >= condition")
	}
}

func TestWhereDistinctOn(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereDistinctOn("name")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "DISTINCT") {
		t.Error("WhereDistinctOn() should add DISTINCT ON")
	}
}

func TestWhereJsonbEqual(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereJsonbEqual("metadata", "status", "active")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "metadata") || !strings.Contains(sql, "->>") {
		t.Error("WhereJsonbEqual() should add JSONB condition")
	}
}

func TestWhereJsonbPathEqual(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereJsonbPathEqual("metadata", []string{"user", "name"}, "John")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "metadata") {
		t.Error("WhereJsonbPathEqual() should add JSONB path condition")
	}
}

func TestWhereJsonbObjectsArrayKeyValueEqual(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereJsonbObjectsArrayKeyValueEqual("tags", "items", "id", "123")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "tags") || !strings.Contains(sql, "@>") {
		t.Error("WhereJsonbObjectsArrayKeyValueEqual() should add JSONB array condition")
	}
}

func TestWhereJsonbPathObjectsArrayKeyValueEqual(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	query := db.NewSelect().Model((*testModel)(nil))
	selector := WhereJsonbPathObjectsArrayKeyValueEqual("metadata", []string{"tags"}, "id", "123")

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "metadata") || !strings.Contains(sql, "@>") {
		t.Error("WhereJsonbPathObjectsArrayKeyValueEqual() should add JSONB path array condition")
	}
}

func TestJsonbPathExpression(t *testing.T) {
	tests := []struct {
		name string
		path []string
		text bool
		want string
	}{
		{
			name: "single path element with text",
			path: []string{"user"},
			text: true,
			want: "?TableAlias.? ->> 'user'",
		},
		{
			name: "multiple path elements with text",
			path: []string{"user", "profile", "name"},
			text: true,
			want: "?TableAlias.? -> 'user' -> 'profile' ->> 'name'",
		},
		{
			name: "multiple path elements without text",
			path: []string{"user", "profile"},
			text: false,
			want: "?TableAlias.? -> 'user' -> 'profile'",
		},
		{
			name: "empty string in path",
			path: []string{"user", "", "name"},
			text: true,
			want: "?TableAlias.? -> 'user' ->> 'name'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jsonbPathExpression(tt.path, tt.text)
			if got != tt.want {
				t.Errorf("jsonbPathExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEscapeJsonPathSegment(t *testing.T) {
	tests := []struct {
		name    string
		segment string
		want    string
	}{
		{
			name:    "no special chars",
			segment: "user",
			want:    "user",
		},
		{
			name:    "single quote",
			segment: "user's",
			want:    "user''s",
		},
		{
			name:    "multiple quotes",
			segment: "it's user's",
			want:    "it''s user''s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeJsonPathSegment(tt.segment)
			if got != tt.want {
				t.Errorf("escapeJsonPathSegment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNestedWhere(t *testing.T) {
	db := newTestDB()
	defer db.Close()

	where := &Where{
		IDs: []string{"1", "2", "3"},
	}

	query := db.NewSelect().Model((*testModel)(nil))
	selector := NestedWhere(*where)

	result := selector(query)
	sql := result.String()

	if !strings.Contains(sql, "id") || !strings.Contains(sql, "IN") {
		t.Error("NestedWhere() should apply Where struct conditions")
	}
}
