package tablesmethods

import (
	"database/sql"
	"fmt"
	"log"
)

func SaveHelper(db *sql.DB, env string, query string, structFields ...any) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	_, err = stmt.Exec(structFields...)
	if err != nil {
		log.Printf("%s: unmatched arguments to insert, err: %v", env, err)
		return fmt.Errorf("%s: unmatched arguments to insert, err: %w", env, err)
	}

	return nil
}

func DeleteByIDHelper(db *sql.DB, env string, query string, id int64) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("%s: failed to prepare the stmt, err: %v", env, err)
		return fmt.Errorf("%s: failed to prepare the stmt, err: %w", env, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("%s: failed to exec the stmt, err: %v", env, err)
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}
