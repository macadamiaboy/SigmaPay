package players

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type Player struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	TgLink     string `json:"tg_link"`
	IsSigma    bool   `json:"is_sigma"`
	PositionID int64  `json:"position_id"`
}

type PaymentCard struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name,omitempty"`
	Surname   string    `json:"surname,omitempty"`
	Price     int       `json:"price"`
	Payed     bool      `json:"payed"`
	EventType string    `json:"event_type,omitempty"`
	EventDate time.Time `json:"event_date,omitempty"`
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

func (p *Player) GetAllSigma(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.players.GetAll"

	rows, err := db.Query("SELECT id, name, surname, tg_link, is_sigma, position_id FROM players WHERE is_sigma = true;")
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

func GetAllPlayersPayments(db *sql.DB, p *Player) (*[]any, error) {
	env := "postgres.tables-methods.players.GetAllPlayersPayments"

	rows, err := db.Query(`
	SELECT pay.id, pl.name, pl.surname, pay.price, pay.payed, et.type, ev.datetime
	FROM players pl
	JOIN payments pay ON pl.id = pay.player_id
	JOIN events ev ON pay.event_id = ev.id
	JOIN pricelist et ON ev.type_id = et.id
	WHERE pl.id = $1;`, p.Id)
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var pCard PaymentCard
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

func GetAllPlayersDebts(db *sql.DB, p *Player) (*[]any, error) {
	env := "postgres.tables-methods.players.GetAllPlayersDebts"

	rows, err := db.Query(`
	SELECT pay.id, pl.name, pl.surname, pay.price, et.type, ev.datetime
	FROM players pl
	JOIN payments pay ON pl.id = pay.player_id
	JOIN events ev ON pay.event_id = ev.id
	JOIN pricelist et ON ev.type_id = et.id
	WHERE pay.payed = false AND pl.id = $1;`, p.Id)
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var pCard PaymentCard
		if err := rows.Scan(&pCard.Id, &pCard.Name, &pCard.Surname, &pCard.Price, &pCard.EventType, &pCard.EventDate); err != nil {
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

func GetTotalDebt(db *sql.DB, p *Player) (*[]any, error) {
	env := "postgres.tables-methods.players.GetTotalDebt"

	rows, err := db.Query(`
	SELECT SUM(pay.price)
	FROM players pl
	JOIN payments pay ON pl.id = pay.player_id
	WHERE pay.payed = false AND pl.id = $1;`, p.Id)
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var pCard PaymentCard
		if err := rows.Scan(&pCard.Price); err != nil {
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
