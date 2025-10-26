package bunhelpers

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

func TestIsConstraintError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "constraint error",
			err:  errors.New("duplicate key value violates unique constraint"),
			want: true,
		},
		{
			name: "constraint error with extra context",
			err:  errors.New("ERROR: duplicate key value violates unique constraint \"users_email_key\""),
			want: true,
		},
		{
			name: "not a constraint error",
			err:  errors.New("some other error"),
			want: false,
		},
		{
			name: "connection error",
			err:  errors.New("connection refused"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsConstraintError(tt.err); got != tt.want {
				t.Errorf("IsConstraintError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "sql.ErrNoRows",
			err:  sql.ErrNoRows,
			want: true,
		},
		{
			name: "wrapped sql.ErrNoRows",
			err:  fmt.Errorf("failed to find user: %w", sql.ErrNoRows),
			want: true,
		},
		{
			name: "other error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFoundError(tt.err); got != tt.want {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.want)
			}
		})
	}
}
