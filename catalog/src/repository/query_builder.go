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

func QueryBuilderOrDefault(runner squirrel.BaseRunner, fallback *squirrel.StatementBuilderType) *squirrel.StatementBuilderType {
	if runner != nil {
		return QueryBuilder(runner)
	}

	return fallback
}
