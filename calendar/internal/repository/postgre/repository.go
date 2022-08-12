package postgre

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *sql.DB
}

func NewRepository(dsn string) *Repository {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return &Repository{
		db: db,
	}
}
