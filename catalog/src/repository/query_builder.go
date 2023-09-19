package repository

import (
	"github.com/Masterminds/squirrel"
)

func QueryBuilder(runner squirrel.BaseRunner) *squirrel.StatementBuilderType {
	qb := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(runner)
	return &qb
}
