package tablesmethods

import (
	"database/sql"
	"fmt"
	"log"
)

type EventType struct {
	Id    int64
	Type  string
	Price int
}

func (e *EventType) Save(db *sql.DB) error {
	env := "postgres.tables-methods.pricelist.Save"
	query := "INSERT INTO pricelist(type, price) VALUES($1, $2);"

	return SaveHelper(db, env, query, e.Type, e.Price)
}

func (e *EventType) GetByID(db *sql.DB, id int64) (*EventType, error) {
	env := "postgres.tables-methods.pricelist.GetByID"

	stmt, err := db.Prepare("SELECT * FROM pricelist WHERE id = $1")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfEvent int64
	var typeOfEvent string
	var priceOfEvent int
	err = stmt.QueryRow(id).Scan(&idOfEvent, &typeOfEvent, &priceOfEvent)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res EventType = EventType{Id: idOfEvent, Type: typeOfEvent, Price: priceOfEvent}

	return &res, nil
}

func (e *EventType) DeleteByID(db *sql.DB, id int64) error {
	env := "postgres.tables-methods.pricelist.DeleteByID"
	query := "DELETE FROM pricelist WHERE id = $1;"

	return DeleteByIDHelper(db, env, query, id)
}
