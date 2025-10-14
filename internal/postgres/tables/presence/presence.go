package presence

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type PlayerPresence struct {
	Id       int64
	EventID  int64
	PlayerID int64
}

func (p *PlayerPresence) Save(db *sql.DB) error {
	env := "postgres.tables-methods.presence.Save"
	query := "INSERT INTO player_presence(event_id, player_id) VALUES($1, $2);"

	return tablesmethods.ExecHelper(db, env, query, p.EventID, p.PlayerID)
}

func (p *PlayerPresence) Update(db *sql.DB) error {
	env := "postgres.tables-methods.presence.Update"
	query := "UPDATE events SET event_id = $2, player_id = $3 WHERE id = $1;"

	return tablesmethods.ExecHelper(db, env, query, p.Id, p.EventID, p.PlayerID)
}

func (p *PlayerPresence) Get(db *sql.DB) (*PlayerPresence, error) {
	env := "postgres.tables-methods.presence.GetByID"

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

func (p *PlayerPresence) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.presence.Delete"
	query := "DELETE FROM player_presence WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, p.Id)
}
