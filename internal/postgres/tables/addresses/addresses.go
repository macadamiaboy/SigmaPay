package addresses

import (
	"database/sql"
	"fmt"
	"log"

	tablesmethods "github.com/macadamiaboy/SigmaPay/internal/postgres/tables"
)

type Address struct {
	Id       int64
	Street   string
	House    int
	Building int
}

func (a *Address) Save(db *sql.DB) error {
	env := "postgres.tables-methods.addresses.Save"
	query := "INSERT INTO addresses(street, house, building) VALUES($1, $2, $3);"

	return tablesmethods.ExecHelper(db, env, query, a.Street, a.House, a.Building)
}

func (a *Address) Update(db *sql.DB) error {
	env := "postgres.tables-methods.addresses.Update"
	query := "UPDATE addresses SET street = $2, house = $3, building = $4 WHERE id = $1;"

	record, err := a.Get(db)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	address, ok := record.(*Address)
	if ok {
		if a.Street == "" {
			a.Street = address.Street
		}
		if a.House == 0 {
			a.House = address.House
		}
		if a.Building == 0 {
			a.Building = address.Building
		}
	}

	return tablesmethods.ExecHelper(db, env, query, a.Id, a.Street, a.House, a.Building)
}

func (a *Address) Get(db *sql.DB) (any, error) {
	env := "postgres.tables-methods.addresses.Get"

	stmt, err := db.Prepare("SELECT * FROM addresses WHERE id = $1;")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfAddress int64
	var streetOfAddress string
	var houseOfAddress int
	var buildingOfAddress int
	err = stmt.QueryRow(a.Id).Scan(&idOfAddress, &streetOfAddress, &houseOfAddress, &buildingOfAddress)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Address = Address{Id: idOfAddress, Street: streetOfAddress, House: houseOfAddress, Building: buildingOfAddress}

	return &res, nil
}

func (p *Address) GetAll(db *sql.DB) (*[]any, error) {
	env := "postgres.tables-methods.addresses.GetAll"

	rows, err := db.Query("SELECT id, street, house, building FROM addresses;")
	if err != nil {
		log.Printf("%s: failed to execute the query, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to execute the query, err: %w", env, err)
	}
	defer rows.Close()

	var collection []any
	for rows.Next() {
		var address Address
		if err := rows.Scan(&address.Id, &address.Street, &address.House, &address.Building); err != nil {
			log.Printf("%s: failed to get the payment, err: %v", env, err)
			return nil, fmt.Errorf("%s: failed to get the payment, err: %w", env, err)
		}
		collection = append(collection, address)
	}

	if err = rows.Err(); err != nil {
		log.Printf("%s: error occured with table rows, err: %v", env, err)
		return nil, fmt.Errorf("%s: error occured with table rows, err: %w", env, err)
	}

	return &collection, nil
}

func (a *Address) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.addresses.Delete"
	query := "DELETE FROM addresses WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, a.Id)
}
