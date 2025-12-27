package events

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
	"github.com/macadamiaboy/SigmaPay/internal/postgres/tables/players"
)

type Event struct {
	Id        int64     `json:"id"`
	TypeID    int64     `json:"type_id"`
	AddressID int64     `json:"address_id"`
	DateTime  time.Time `json:"datetime"`
}

type EventCard struct {
	Id       int64     `json:"id"`
	Type     string    `json:"type"`
	Address  string    `json:"address"`
	DateTime time.Time `json:"datetime"`
}

func (e *Event) Save(db *sql.DB) error {
	env := "postgres.tables-methods.events.Save"
	query := "INSERT INTO events(type_id, address_id, datetime) VALUES($1, $2, $3);"

	return tablesmethods.ExecHelper(db, env, query, e.TypeID, e.AddressID, e.DateTime)
}

func (e *Event) Update(db *sql.DB) error {
	env := "postgres.tables-methods.events.Update"
	query := "UPDATE events SET type_id = $2, address_id = $3, datetime = $4 WHERE id = $1;"

	record, err := e.Get(db)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	event, ok := record.(*Event)
	if ok {
		if e.TypeID == 0 {
			e.TypeID = event.TypeID
		}
		if e.AddressID == 0 {
			e.AddressID = event.AddressID
		}
		if e.DateTime.IsZero() {
			e.DateTime = event.DateTime
		}
	}

	return tablesmethods.ExecHelper(db, env, query, e.Id, e.TypeID, e.AddressID, e.DateTime)
}

func (e *Event) Get(db *sql.DB) (any, error) {
	env := "postgres.tables-methods.events.Get"

	stmt, err := db.Prepare("SELECT * FROM events WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfEvent int64
	var idOfType int64
	var idOfAddress int64
	var timeOfEvent time.Time
	err = stmt.QueryRow(e.Id).Scan(&idOfEvent, &idOfType, &idOfAddress, &timeOfEvent)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Event = Event{Id: idOfEvent, TypeID: idOfType, AddressID: idOfAddress, DateTime: timeOfEvent}

	return &res, nil
}

func (e *Event) GetAll(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.events.GetAll"

	rows, err := db.Query("SELECT id, type_id, address_id, datetime FROM events;")
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.Id, &event.TypeID, &event.AddressID, &event.DateTime); err != nil {
			log.Printf("%s: failed to get the event, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the event, err: %w", env, err)
		}
		collection = append(collection, event)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (e *Event) GetAllByMonth(db *sql.DB, reqMonth string, reqYear string) (*[]any, error) {
	env := "postgres.tables-methods.events.GetAllByMonth"

	rows, err := db.Query(`
	SELECT e.id, p.type, CONCAT(a.street, ' ', a.house, 'к', a.building) AS address, e.datetime
	FROM events e
	JOIN pricelist p ON e.type_id = p.id
	JOIN addresses a ON e.address_id = a.id
	WHERE EXTRACT(MONTH FROM e.datetime) = $1 AND EXTRACT(YEAR FROM e.datetime) = $2;`, reqMonth, reqYear)
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var eventCard EventCard
		if err := rows.Scan(&eventCard.Id, &eventCard.Type, &eventCard.Address, &eventCard.DateTime); err != nil {
			log.Printf("%s: failed to get the event, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the event, err: %w", env, err)
		}
		collection = append(collection, eventCard)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (e *Event) getAllByType(db *sql.DB, env string, reqType string) (*[]any, error) {
	rows, err := db.Query(`
	SELECT EXTRACT(MONTH FROM e.datetime)
	FROM events e
	JOIN pricelist p
	ON e.type_id = p.id
	WHERE p.type = $1;`, reqType)
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var month string
		if err := rows.Scan(&month); err != nil {
			log.Printf("%s: failed to get the event, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the event, err: %w", env, err)
		}
		collection = append(collection, month)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (e *Event) GetAllTrainings(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.events.GetAllTrainings"
	reqType := "Тренировка"

	return e.getAllByType(db, env, reqType)
}

func (e *Event) GetAllGames(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.events.GetAllGames"
	reqType := "Игра"

	return e.getAllByType(db, env, reqType)
}

func (e *Event) GetAllEventPayments(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.events.GetAllEventPayments"

	rows, err := db.Query(`
	SELECT pay.id, pl.name, pl.surname, pay.price, pay.payed, et.type, ev.datetime
	FROM events ev
	JOIN payments pay ON ev.id = pay.event_id
	JOIN players pl ON pay.player_id = pl.id
	JOIN pricelist et ON ev.type_id = et.id
	WHERE ev.id = $1;`, e.Id)
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var pCard players.PaymentCard
		if err := rows.Scan(&pCard.Id, &pCard.Name, &pCard.Surname, &pCard.Price, &pCard.Payed, &pCard.EventType, &pCard.EventDate); err != nil {
			log.Printf("%s: failed to get the player, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the player, err: %w", env, err)
		}
		collection = append(collection, pCard)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (e *Event) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.events.Delete"
	query := "DELETE FROM events WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, e.Id)
}
