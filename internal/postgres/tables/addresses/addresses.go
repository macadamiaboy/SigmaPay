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

	return tablesmethods.SaveHelper(db, env, query, a.Street, a.House, a.Building)
}

func GetByID(db *sql.DB, id int64) (*Address, error) {
	env := "postgres.tables-methods.addresses.GetByID"

	stmt, err := db.Prepare("SELECT * FROM addresses WHERE id = $1")
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return nil, fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	var idOfAddress int64
	var streetOfAddress string
	var houseOfAddress int
	var buildingOfAddress int
	err = stmt.QueryRow(id).Scan(&idOfAddress, &streetOfAddress, &houseOfAddress, &buildingOfAddress)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	var res Address = Address{Id: idOfAddress, Street: streetOfAddress, House: houseOfAddress, Building: buildingOfAddress}

	return &res, nil
}

func DeleteByID(db *sql.DB, id int64) error {
	env := "postgres.tables-methods.addresses.DeleteByID"
	query := "DELETE FROM addresses WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, id)
}

func (a *Address) Delete(db *sql.DB) error {
	env := "postgres.tables-methods.addresses.Delete"
	query := "DELETE FROM addresses WHERE id = $1;"

	return tablesmethods.DeleteByIDHelper(db, env, query, a.Id)
}
