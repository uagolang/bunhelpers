package bunhelpers

import (
	"database/sql"
	"errors"
	"strings"
)

func IsConstraintError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
