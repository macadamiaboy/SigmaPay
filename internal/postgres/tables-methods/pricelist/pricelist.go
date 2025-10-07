package pricelist

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables-methods"
)

type EventType struct {
	Id    int64
	Type  string
	Price int
}

func (e *EventType) Save(db *sql.DB) error {
	env := "postgres.tables-methods.pricelist.Save"
	query := "INSERT INTO pricelist(type, price) VALUES($1, $2);"

	return tablesmethods.SaveHelper(db, env, query, e.Type, e.Price)
}

func GetByID(db *sql.DB, id int64) (*EventType, error) {
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

func DeleteByID(db *sql.DB, id int64) error {
	env := "postgres.tables-methods.pricelist.DeleteByID"
	query := "DELETE FROM pricelist WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, id)
}

func (e *EventType) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.pricelist.DeleteByID"
	query := "DELETE FROM pricelist WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, e.Id)
}
