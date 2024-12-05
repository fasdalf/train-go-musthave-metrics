package catchable

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsPgConnectionError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr); pgErr != nil && pgerrcode.IsConnectionException(pgErr.Code) {
		return true
	}
	return false
}
