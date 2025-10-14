package payments

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type Payment struct {
	Id       int64
	PlayerID int64
	Price    int
	Payed    bool
}

func (p *Payment) Save(db *sql.DB) error {
	env := "postgres.tables-methods.payments.Save"
	query := "INSERT INTO payments(player_id, price, payed) VALUES($1, $2, $3);"

	return tablesmethods.ExecHelper(db, env, query, p.PlayerID, p.Price, p.Payed)
}

func (p *Payment) Update(db *sql.DB) error {
	env := "postgres.tables-methods.payments.Update"
	query := "UPDATE payments SET player_id = $2, price = $3, payed = $4 WHERE id = $1;"

	return tablesmethods.ExecHelper(db, env, query, p.Id, p.PlayerID, p.Price, p.Payed)
}

func GetByID(db *sql.DB, id int64) (*Payment, error) {
	env := "postgres.tables-methods.payments.GetByID"

	stmt, err := db.Prepare("SELECT * FROM payments WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfPayment int64
	var idOfPlayer int64
	var price int
	var payed bool
	err = stmt.QueryRow(id).Scan(&idOfPayment, &idOfPlayer, &price, &payed)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Payment = Payment{Id: idOfPayment, PlayerID: idOfPlayer, Price: price, Payed: payed}

	return &res, nil
}

func (p *Payment) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.payments.Delete"
	query := "DELETE FROM payments WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, p.Id)
}
