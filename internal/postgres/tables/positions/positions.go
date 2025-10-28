package positions

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type Position struct {
	Id       int64  `json:"id"`
	Position string `json:"position"`
}

func (p *Position) Save(db *sql.DB) error {
	env := "postgres.tables-methods.positions.Save"
	query := "INSERT INTO positions(position) VALUES($1);"

	return tablesmethods.ExecHelper(db, env, query, p.Position)
}

func (p *Position) Update(db *sql.DB) error {
	env := "postgres.tables-methods.positions.Update"
	query := "UPDATE events SET position = $2 WHERE id = $1;"

	return tablesmethods.ExecHelper(db, env, query, p.Id, p.Position)
}

func (p *Position) Get(db *sql.DB) (any, error) {
	env := "postgres.tables-methods.positions.Get"
	stmt, err := db.Prepare("SELECT * FROM positions WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfPosition int64
	var position string
	err = stmt.QueryRow(p.Id).Scan(&idOfPosition, &position)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Position = Position{Id: idOfPosition, Position: position}

	return &res, nil
}

func (p *Position) GetAll(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.positions.GetAll"

	rows, err := db.Query("SELECT id, position FROM positions;")
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var position Position
		if err := rows.Scan(&position.Id, &position.Position); err != nil {
			log.Printf("%s: failed to get the position, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the position, err: %w", env, err)
		}
		collection = append(collection, position)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (p *Position) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.positions.Delete"
	query := "DELETE FROM positions WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, p.Id)
}
