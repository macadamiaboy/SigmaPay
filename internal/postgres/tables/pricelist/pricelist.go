package pricelist

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type EventType struct {
	Id    int64  `json:"id"`
	Type  string `json:"type"`
	Price int    `json:"price"`
}

func (e *EventType) Save(db *sql.DB) error {
	env := "postgres.tables-methods.pricelist.Save"
	query := "INSERT INTO pricelist(type, price) VALUES($1, $2);"

	return tablesmethods.ExecHelper(db, env, query, e.Type, e.Price)
}

func (e *EventType) Update(db *sql.DB) error {
	env := "postgres.tables-methods.pricelist.Update"
	query := "UPDATE pricelist SET type = $2, price = $3 WHERE id = $1;"

	return tablesmethods.ExecHelper(db, env, query, e.Id, e.Type, e.Price)
}

func GetByID(db *sql.DB, id int64) (*EventType, error) {
	env := "postgres.tables-methods.pricelist.GetByID"

	stmt, err := db.Prepare("SELECT * FROM pricelist WHERE id = $1;")
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

func GetAll(db *sql.DB) (*[]EventType, error) {
	env := "postgres.tables-methods.pricelist.GetAll"

	rows, err := db.Query("SELECT id, type, price FROM pricelist;")
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []EventType
	for rows.Next() {
		var eventType EventType
		if err := rows.Scan(&eventType.Id, &eventType.Type, &eventType.Price); err != nil {
			log.Printf("%s: failed to get the eventType, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the eventType, err: %w", env, err)
		}
		collection = append(collection, eventType)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (e *EventType) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.pricelist.Delete"
	query := "DELETE FROM pricelist WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, e.Id)
}
