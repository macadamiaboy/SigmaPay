package payments

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type Payment struct {
	Id       int64 `json:"id"`
	PlayerID int64 `json:"player_id"`
	EventID  int64 `json:"event_id"`
	Price    int   `json:"price"`
	Payed    bool  `json:"payed"`
}

func (p *Payment) Save(db *sql.DB) error {
	env := "postgres.tables-methods.payments.Save"
	query := "INSERT INTO payments(player_id, event_id, price, payed) VALUES($1, $2, $3, $4);"

	return tablesmethods.ExecHelper(db, env, query, p.PlayerID, p.EventID, p.Price, p.Payed)
}

func (p *Payment) Update(db *sql.DB) error {
	env := "postgres.tables-methods.payments.Update"
	query := "UPDATE payments SET player_id = $2, event_id = $3, price = $4, payed = $5 WHERE id = $1;"

	record, err := p.Get(db)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	payment, ok := record.(*Payment)
	if ok {
		if p.PlayerID == 0 {
			p.PlayerID = payment.PlayerID
		}
		if p.EventID == 0 {
			p.EventID = payment.EventID
		}
		if p.Price == 0 {
			p.Price = payment.Price
		}
		//payed field is not touched
	}

	return tablesmethods.ExecHelper(db, env, query, p.Id, p.PlayerID, p.EventID, p.Price, p.Payed)
}

func (p *Payment) Get(db *sql.DB) (any, error) {
	env := "postgres.tables-methods.payments.Get"

	stmt, err := db.Prepare("SELECT * FROM payments WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfPayment int64
	var idOfPlayer int64
	var idOfEvent int64
	var price int
	var payed bool
	err = stmt.QueryRow(p.Id).Scan(&idOfPayment, &idOfPlayer, &idOfEvent, &price, &payed)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Payment = Payment{Id: idOfPayment, PlayerID: idOfPlayer, EventID: idOfEvent, Price: price, Payed: payed}

	return &res, nil
}

func (p *Payment) GetAll(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.payments.GetAll"

	rows, err := db.Query("SELECT id, player_id, event_id, price, payed FROM payments;")
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var payment Payment
		if err := rows.Scan(&payment.Id, &payment.PlayerID, &payment.EventID, &payment.Price, &payment.Payed); err != nil {
			log.Printf("%s: failed to get the payment, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the payment, err: %w", env, err)
		}
		collection = append(collection, payment)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (p *Payment) GetNotPayed(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.payments.GetNotPayed"

	rows, err := db.Query("SELECT id, player_id, event_id, price, payed FROM payments WHERE payed = false;")
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var payment Payment
		if err := rows.Scan(&payment.Id, &payment.PlayerID, &payment.EventID, &payment.Price, &payment.Payed); err != nil {
			log.Printf("%s: failed to get the payment, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the payment, err: %w", env, err)
		}
		collection = append(collection, payment)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (p *Payment) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.payments.Delete"
	query := "DELETE FROM payments WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, p.Id)
}
