package players

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type Player struct {
	Id         int64
	Name       string
	Surname    string
	TgLink     string
	IsSigma    bool
	PositionID int64
}

func (p *Player) Save(db *sql.DB) error {
	env := "postgres.tables-methods.players.Save"
	query := "INSERT INTO players(name, surname, tg_link, is_sigma, position_id) VALUES($1, $2, $3, $4, $5);"

	return tablesmethods.ExecHelper(db, env, query, p.Name, p.Surname, p.TgLink, p.IsSigma, p.PositionID)
}

func (p *Player) Update(db *sql.DB) error {
	env := "postgres.tables-methods.players.Update"
	query := "UPDATE events SET name = $2, surname = $3, tg_link = $4, is_sigma = $5, position_id = $6 WHERE id = $1;"

	record, err := p.Get(db)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	player, ok := record.(*Player)
	if ok {
		if p.Name == "" {
			p.Name = player.Name
		}
		if p.Surname == "" {
			p.Surname = player.Surname
		}
		if p.TgLink == "" {
			p.TgLink = player.TgLink
		}
		if p.PositionID == 0 {
			p.PositionID = player.PositionID
		}
		//we skip isSigma check
	}

	return tablesmethods.ExecHelper(db, env, query, p.Id, p.Name, p.Surname, p.TgLink, p.IsSigma, p.PositionID)
}

func (p *Player) Get(db *sql.DB) (any, error) {
	env := "postgres.tables-methods.players.Get"

	stmt, err := db.Prepare("SELECT * FROM players WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfPlayer int64
	var nameOfPlayer string
	var surnameOfPlayer string
	var tgLink string
	var isSigma bool
	var idOfPosition int64
	err = stmt.QueryRow(p.Id).Scan(&idOfPlayer, &nameOfPlayer, &surnameOfPlayer, &tgLink, &isSigma, &idOfPosition)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Player = Player{Id: idOfPlayer, Name: nameOfPlayer, Surname: surnameOfPlayer, TgLink: tgLink, IsSigma: isSigma, PositionID: idOfPosition}

	return &res, nil
}

func (p *Player) GetAll(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.players.GetAll"

	rows, err := db.Query("SELECT id, name, surname, tg_link, is_sigma, position_id FROM players;")
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var player Player
		if err := rows.Scan(&player.Id, &player.Name, &player.Surname, &player.TgLink, &player.IsSigma, &player.PositionID); err != nil {
			log.Printf("%s: failed to get the player, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the player, err: %w", env, err)
		}
		collection = append(collection, player)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (p *Player) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.players.Delete"
	query := "DELETE FROM players WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, p.Id)
}
