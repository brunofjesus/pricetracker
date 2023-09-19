package model

const StoreTableName = "store"

type Store struct {
	StoreId int64  `db:"store_id"`
	Slug    string `db:"slug"`
	Name    string `db:"name"`
	Website string `db:"website"`
	Active  bool   `db:"active"`
}
