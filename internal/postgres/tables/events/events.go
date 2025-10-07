package events

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type Event struct {
	Id        int64
	TypeID    int64
	AddressID int64
	Time      time.Time
}

func (e *Event) Save(db *sql.DB) error {
	env := "postgres.tables-methods.events.Save"
	query := "INSERT INTO events(type_id, address_id, time) VALUES($1, $2, $3);"

	return tablesmethods.SaveHelper(db, env, query, e.TypeID, e.AddressID, e.Time)
}

func GetByID(db *sql.DB, id int64) (*Event, error) {
	env := "postgres.tables-methods.events.GetByID"

	stmt, err := db.Prepare("SELECT * FROM events WHERE id = $1")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfEvent int64
	var idOfType int64
	var idOfAddress int64
	var timeOfEvent time.Time
	err = stmt.QueryRow(id).Scan(&idOfEvent, &idOfType, &idOfAddress, &timeOfEvent)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Event = Event{Id: idOfEvent, TypeID: idOfType, AddressID: idOfAddress, Time: timeOfEvent}

	return &res, nil
}

func DeleteByID(db *sql.DB, id int64) error {
	env := "postgres.tables-methods.events.DeleteByID"
	query := "DELETE FROM events WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, id)
}

func (e *Event) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.events.Delete"
	query := "DELETE FROM events WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, e.Id)
}
