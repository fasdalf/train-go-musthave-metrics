package catchable

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsPgConnectionError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"yes",
			args{err: &pgconn.PgError{Code: "08000"}},
			true,
		},
		{
			"no",
			args{err: fmt.Errorf("some error")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPgConnectionError(tt.args.err); got != tt.want {
				t.Errorf("IsPgConnectionError() = %v, want %v", got, tt.want)
			}
		})
	}
}
