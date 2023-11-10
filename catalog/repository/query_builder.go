package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

func QueryBuilder(runner squirrel.BaseRunner) *squirrel.StatementBuilderType {
	qb := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(runner)
	return &qb
}

func QueryBuilderOrDefault(tx *sql.Tx, fallback *squirrel.StatementBuilderType) *squirrel.StatementBuilderType {
	if tx != nil {
		return QueryBuilder(tx)
	}

	return fallback
}
