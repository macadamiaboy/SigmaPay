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

	return tablesmethods.SaveHelper(db, env, query, p.Name, p.Surname, p.TgLink, p.IsSigma, p.PositionID)
}

func GetByID(db *sql.DB, id int64) (*Player, error) {
	env := "postgres.tables-methods.players.GetByID"

	stmt, err := db.Prepare("SELECT * FROM players WHERE id = $1")
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
	err = stmt.QueryRow(id).Scan(&idOfPlayer, &nameOfPlayer, &surnameOfPlayer, &tgLink, &isSigma, &idOfPosition)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Player = Player{Id: idOfPlayer, Name: nameOfPlayer, Surname: surnameOfPlayer, TgLink: tgLink, IsSigma: isSigma, PositionID: idOfPosition}

	return &res, nil
}

func DeleteByID(db *sql.DB, id int64) error {
	env := "postgres.tables-methods.players.DeleteByID"
	query := "DELETE FROM players WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, id)
}

func (p *Player) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.players.Delete"
	query := "DELETE FROM players WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, p.Id)
}
