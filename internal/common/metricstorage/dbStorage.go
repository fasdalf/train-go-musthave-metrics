package metricstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DBStorage store metrics in DB
// TODO: use context and handle errors in all methods.
type DBStorage struct {
	db *sql.DB
}

type DBBatch struct {
	tx *sql.Tx
}

type execSQL func(query string, args ...any) (sql.Result, error)

func NewDBStorage(db *sql.DB, ctx context.Context) (s *DBStorage, err error) {
	s = &DBStorage{
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
	// если вместо _ = tx.Rollback() в каждом условии здесь вызвать defer tx.Rollback()

	// создаём таблицу целочисленных счетчиков и необходимый индекс
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gauge (
			name varchar(250) NOT NULL,
			value double PRECISION NOT NULL
		)
    `)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS gauge_name_udx ON gauge (name)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// создаём таблицу дробных счетчиков и необходимый индекс
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS counter (
			name varchar(250) NOT NULL,
			value bigint NOT NULL
		)
    `)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS counter_name_udx ON counter (name)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// коммитим транзакцию
	return tx.Commit()
}

func updateCounterSQL(exec execSQL, key string, value int) error {
	_, err := exec(`
        INSERT INTO counter (name, value)
		VALUES (@name, @value)
		ON CONFLICT (name) DO UPDATE SET value = counter.value+EXCLUDED.value;
    `, pgx.NamedArgs{"name": key, "value": value})
	return err

}

func (s *DBStorage) UpdateCounter(key string, value int) error {
	return updateCounterSQL(s.db.Exec, key, value)
}

func updateGaugeSQL(exec execSQL, key string, value float64) error {
	_, err := exec(`
        INSERT INTO gauge (name, value)
		VALUES (@name, @value)
		ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value;
    `, pgx.NamedArgs{"name": key, "value": value})
	return err
}
func (s *DBStorage) UpdateGauge(key string, value float64) error {
	return updateGaugeSQL(s.db.Exec, key, value)
}

func (s *DBStorage) GetCounter(key string) (r int, err error) {
	row := s.db.QueryRow(`
        SELECT t.value FROM counter t WHERE t.name = @name
    `, pgx.NamedArgs{"name": key})
	err = row.Scan(&r)
	if err != nil {
		return 0, err
	}

	return r, err
}

func (s *DBStorage) GetGauge(key string) (r float64, err error) {
	row := s.db.QueryRow(`
        SELECT t.value FROM gauge t WHERE t.name = @name
    `, pgx.NamedArgs{"name": key})
	err = row.Scan(&r)
	if err != nil {
		return 0, err
	}

	return r, err
}

func (s *DBStorage) HasCounter(key string) (r bool, err error) {
	row := s.db.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM counter t WHERE t.name = @name)
    `, pgx.NamedArgs{"name": key})
	err = row.Scan(&r)
	if err != nil {
		return false, err
	}

	return r, err
}

func (s *DBStorage) HasGauge(key string) (r bool, err error) {
	row := s.db.QueryRow(`
        SELECT EXISTS(SELECT 1 FROM gauge t WHERE t.name = @name)
    `, pgx.NamedArgs{"name": key})
	err = row.Scan(&r)
	if err != nil {
		return false, err
	}

	return r, err
}

func (s *DBStorage) ListGauges() ([]string, error) {
	keys := []string{}
	rows, err := s.db.Query(`
        SELECT t.name FROM gauge t
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var k string
		if err = rows.Scan(&s); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *DBStorage) ListCounters() ([]string, error) {
	keys := []string{}
	rows, err := s.db.Query(`
        SELECT t.name FROM counter t
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var k string
		if err = rows.Scan(&s); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *DBStorage) StartBatch(ctx context.Context) (basicBatch, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &DBBatch{tx: tx}, nil
}

func (b *DBBatch) UpdateCounter(key string, value int) error {
	err := updateCounterSQL(b.tx.Exec, key, value)
	if err != nil {
		err = errors.Join(err, b.rollback())
	}
	return err
}
func (b *DBBatch) UpdateGauge(key string, value float64) error {
	err := updateGaugeSQL(b.tx.Exec, key, value)
	if err != nil {
		err = errors.Join(err, b.rollback())
	}
	return err
}
func (b *DBBatch) Commit() error {
	return b.tx.Commit()
}
func (b *DBBatch) rollback() error {
	return b.tx.Commit()
}
