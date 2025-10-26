package bunhelpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Selector func(*bun.SelectQuery) *bun.SelectQuery

// ApplyIf applies the provided Selectors when condition is true.
// Otherwise returns nil Selector func, meaning the provided Selectors will not be applied to the query.
func ApplyIf(cond bool, selectors ...Selector) Selector {
	if !cond {
		return nil
	}
	return Apply(selectors...)
}

// Apply is the same as bun.SelectQuery.Apply.
// It combines multiple Selectors into one.
func Apply(selectors ...Selector) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		for _, selector := range selectors {
			if selector != nil {
				q = selector(q)
			}
		}
		return q
	}
}

// OrGroup adds a group to WHERE clause, prefixed by OR if there are other conditions before it.
func OrGroup(selectors ...Selector) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return Apply(selectors...)(sq)
		})
	}
}

// AndGroup adds a group to WHERE clause, prefixed by AND if there are other conditions before it.
func AndGroup(selectors ...Selector) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return Apply(selectors...)(sq)
		})
	}
}

// Or adds AND group to WHERE clause, in which all conditions are separated by OR.
func Or(selectors ...Selector) Selector {
	return AndGroup(Map(selectors, func(s Selector, _ int) Selector {
		return OrGroup(s)
	})...)
}

// UseWhere allows to reuse the Where.Where() common logic as a Selector.
func UseWhere(where Where) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return where.Where(q)
	}
}

// WhereJsonbEqual compares a JSONB string field to a parameterized value safely using text extraction.
func WhereJsonbEqual(col string, field string, value any) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.?->>? = ?", bun.Ident(col), field, value)
	}
}

// WhereJsonbPathEqual compares a JSONB string field located at the provided path to a value.
// Path elements are applied in order using -> for intermediate levels and ->> for the final key.
func WhereJsonbPathEqual(col string, path []string, value any) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(jsonbPathExpression(path, true)+" = ?", bun.Ident(col), value)
	}
}

// WhereJsonbObjectsArrayKeyValueEqual checks that a JSONB array of objects (at key) contains
// an object where field == value. Built with JSONB functions to avoid string interpolation.
func WhereJsonbObjectsArrayKeyValueEqual(col string, key, field string, value any) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(
			`?TableAlias.? -> ? @> jsonb_build_array(jsonb_build_object(?::text, ?::text))`,
			bun.Ident(col), key, field, value,
		)
	}
}

// WhereJsonbPathObjectsArrayKeyValueEqual checks that the JSONB array of objects located at the provided path
// contains an object where field == value.
func WhereJsonbPathObjectsArrayKeyValueEqual(col string, path []string, field string, value any) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(
			jsonbPathExpression(path, false)+" @> jsonb_build_array(jsonb_build_object(?::text, ?::text))",
			bun.Ident(col), field, value,
		)
	}
}

func jsonbPathExpression(path []string, text bool) string {
	expr := "?TableAlias.?"
	for idx, segment := range path {
		if segment == "" {
			continue
		}
		operator := "->"
		if text && idx == len(path)-1 {
			operator = "->>"
		}
		expr += fmt.Sprintf(" %s '%s'", operator, escapeJsonPathSegment(segment))
	}
	return expr
}

func escapeJsonPathSegment(segment string) string {
	return strings.ReplaceAll(segment, "'", "''")
}

func WhereEqual(col string, value any) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? = ?", bun.Ident(col), value)
	}
}

func WhereNull(col string) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? is null", bun.Ident(col))
	}
}

func WhereNotNull(col string) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? is not null", bun.Ident(col))
	}
}

func WhereDistinctOn(col string) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.DistinctOn(col).OrderExpr(fmt.Sprintf("%s, id", col))
	}
}

func WhereNotEqual(col string, value any) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? != ?", bun.Ident(col), value)
	}
}

func WhereIn(col string, values any) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? IN (?)", bun.Ident(col), bun.In(values))
	}
}

func WhereNotIn(col string, values any) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		q.Table()

		return q.Where("?TableAlias.? NOT IN (?)", bun.Ident(col), bun.In(values))
	}
}

func WhereContains(col string, substr string) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? ILIKE ?", bun.Ident(col), "%"+substr+"%")
	}
}

func WhereBegins(col string, substr string) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? ILIKE ?", bun.Ident(col), substr+"%")
	}
}

func WhereEnds(col string, substr string) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? ILIKE ?", bun.Ident(col), "%"+substr)
	}
}

func WhereBefore(col string, t time.Time) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? <= ?", bun.Ident(col), t)
	}
}

func WhereAfter(col string, t time.Time) Selector {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("?TableAlias.? >= ?", bun.Ident(col), t)
	}
}
