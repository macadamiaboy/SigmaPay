package presence

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type PlayerPresence struct {
	Id       int64 `json:"id"`
	EventID  int64 `json:"event_id"`
	PlayerID int64 `json:"player_id"`
}

func (p *PlayerPresence) Save(db *sql.DB) error {
	env := "postgres.tables-methods.presence.Save"
	query := "INSERT INTO player_presence(event_id, player_id) VALUES($1, $2);"

	return tablesmethods.ExecHelper(db, env, query, p.EventID, p.PlayerID)
}

func (p *PlayerPresence) Update(db *sql.DB) error {
	env := "postgres.tables-methods.presence.Update"
	query := "UPDATE events SET event_id = $2, player_id = $3 WHERE id = $1;"

	record, err := p.Get(db)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	playerPresence, ok := record.(*PlayerPresence)
	if ok {
		if p.EventID == 0 {
			p.EventID = playerPresence.EventID
		}
		if p.PlayerID == 0 {
			p.PlayerID = playerPresence.PlayerID
		}
	}

	return tablesmethods.ExecHelper(db, env, query, p.Id, p.EventID, p.PlayerID)
}

func (p *PlayerPresence) Get(db *sql.DB) (any, error) {
	env := "postgres.tables-methods.presence.Get"

	stmt, err := db.Prepare("SELECT * FROM player_presence WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfPresence int64
	var idOfEvent int64
	var idOfPlayer int64
	err = stmt.QueryRow(p.Id).Scan(&idOfPresence, &idOfEvent, &idOfPlayer)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res PlayerPresence = PlayerPresence{Id: idOfPresence, EventID: idOfEvent, PlayerID: idOfPlayer}

	return &res, nil
}

func (p *PlayerPresence) GetAll(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.presence.GetAll"

	rows, err := db.Query("SELECT id, event_id, player_id FROM player_presence;")
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var playerPresence PlayerPresence
		if err := rows.Scan(&playerPresence.Id, &playerPresence.EventID, &playerPresence.PlayerID); err != nil {
			log.Printf("%s: failed to get the presence, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the presence, err: %w", env, err)
		}
		collection = append(collection, playerPresence)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (p *PlayerPresence) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.presence.Delete"
	query := "DELETE FROM player_presence WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, p.Id)
}
