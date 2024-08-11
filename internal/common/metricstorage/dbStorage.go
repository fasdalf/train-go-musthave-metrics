package metricstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
)

// DBStorage store metrics in DB
// TODO: use context and handle errors in all methods.
type DBStorage struct {
	r  Retryer
	db *sql.DB
}

type Retryer interface {
	Try(do func() error, isRetryable func(err error) bool) (int, error)
}

func NewDBStorage(db *sql.DB, ctx context.Context, retryer Retryer) (s *DBStorage, err error) {
	s = &DBStorage{
		r:  retryer,
		db: db,
	}

	// IRL it should be done in main() with separate command line flag.
	if err = s.Bootstrap(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (s *DBStorage) Bootstrap(ctx context.Context) error {
	// запускаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// в случае неуспешного коммита все изменения транзакции будут отменены
	defer tx.Rollback()

	// создаём таблицу целочисленных счетчиков и необходимый индекс
	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gauge (
			name varchar(250) NOT NULL,
			value double PRECISION NOT NULL
		)
    `)
	tx.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS gauge_name_udx ON gauge (name)`)

	// создаём таблицу дробных счетчиков и необходимый индекс
	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS counter (
			name varchar(250) NOT NULL,
			value bigint NOT NULL
		)
    `)
	tx.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS counter_name_udx ON counter (name)`)

	// коммитим транзакцию
	return tx.Commit()
}

func (s *DBStorage) UpdateCounter(key string, value int) {
	doJob := func() error {
		_, err := s.db.Exec(`
			INSERT INTO counter (name, value)
			VALUES (@name, @value)
			ON CONFLICT (name) DO UPDATE SET value = counter.value+EXCLUDED.value;
		`, pgx.NamedArgs{"name": key, "value": value})
		return err
	}
	if _, err := s.r.Try(doJob, isPgConnectionError); err != nil {
		slog.Error("UpdateCounter failed", "error", err)
	}
}

func (s *DBStorage) UpdateGauge(key string, value float64) {
	doJob := func() error {
		_, err := s.db.Exec(`
			INSERT INTO gauge (name, value)
			VALUES (@name, @value)
			ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;
		`, pgx.NamedArgs{"name": key, "value": value})
		return err
	}
	if _, err := s.r.Try(doJob, isPgConnectionError); err != nil {
		slog.Error("UpdateGauge failed", "error", err)
	}
}

func (s *DBStorage) GetCounter(key string) (r int) {
	row := s.db.QueryRow(`
        SELECT t.value FROM counter t WHERE t.name = @name
    `, pgx.NamedArgs{"name": key})
	err := row.Scan(&r)
	if err != nil {
		slog.Error("GetCounter failed", "error", err)
		return 0
	}

	return r
}

func (s *DBStorage) GetGauge(key string) (r float64) {
	row := s.db.QueryRow(`
        SELECT t.value FROM gauge t WHERE t.name = @name
    `, pgx.NamedArgs{"name": key})
	err := row.Scan(&r)
	if err != nil {
		slog.Error("GetCounter failed", "error", err)
		return 0
	}

	return r
}

func (s *DBStorage) HasCounter(key string) (r bool) {
	row := s.db.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM counter t WHERE t.name = @name)
    `, pgx.NamedArgs{"name": key})
	err := row.Scan(&r)
	if err != nil {
		slog.Error("HasCounter failed", "error", err)
		return false
	}

	return r
}

func (s *DBStorage) HasGauge(key string) (r bool) {
	row := s.db.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM gauge t WHERE t.name = @name)
    `, pgx.NamedArgs{"name": key})
	err := row.Scan(&r)
	if err != nil {
		slog.Error("HasGauge failed", "error", err)
		return false
	}

	return r
}

func (s *DBStorage) ListGauges() []string {
	keys := []string{}
	rows, err := s.db.Query(`
        SELECT t.name FROM gauge t
    `)
	if err != nil {
		slog.Error("ListGauges failed", "error", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var k string
		if err = rows.Scan(&s); err != nil {
			slog.Error("ListGauges failed", "error", err)
			return nil
		}
		keys = append(keys, k)
	}

	if err = rows.Err(); err != nil {
		slog.Error("ListGauges failed", "error", err)
		return nil
	}

	return keys
}

func (s *DBStorage) ListCounters() []string {
	keys := []string{}
	rows, err := s.db.Query(`
        SELECT t.name FROM counter t
    `)
	if err != nil {
		slog.Error("ListCounters failed", "error", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var k string
		if err = rows.Scan(&s); err != nil {
			slog.Error("ListCounters failed", "error", err)
			return nil
		}
		keys = append(keys, k)
	}

	if err = rows.Err(); err != nil {
		slog.Error("ListCounters failed", "error", err)
		return nil
	}

	return keys
}

func isPgConnectionError(err error) bool {
	pgErr := (*pgconn.PgError)(nil)
	if errors.As(err, &pgErr); pgErr != nil && pgerrcode.IsConnectionException(pgErr.Code) {
		return true
	}
	return false
}
