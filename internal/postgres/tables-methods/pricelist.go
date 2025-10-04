package tablesmethods

import (
	"database/sql"
)

type EventType struct {
	Type  string
	Price int
}

func (e *EventType) Save(db *sql.DB) error {
	env := "postgres.tables-methods.pricelist.Save"
	query := "INSERT INTO pricelist(type, price) VALUES($1, $2);"

	return SaveHelper(db, env, query, e.Type, e.Price)
}
