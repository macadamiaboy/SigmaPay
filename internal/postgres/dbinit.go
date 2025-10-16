package postgres

import (
	"database/sql"
	"fmt"
)

func InitDatabase(db *sql.DB) error {
	if err := initPricelistTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initAddressessTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initPositionsTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initEventsTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initPlayersTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initPlayerPresenceTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	if err := initPaymentsTable(db); err != nil {
		return fmt.Errorf("error occured during init process: %w", err)
	}

	return nil
}

func initPricelistTable(db *sql.DB) error {
	env := "dbinit.initPricelistTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS pricelist(
	    id BIGSERIAL PRIMARY KEY,
	    type VARCHAR(20) NOT NULL,
	    price INTEGER NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	if err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_type ON pricelist(type);"); err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initAddressessTable(db *sql.DB) error {
	env := "dbinit.initAddressessTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS addresses(
	    id BIGSERIAL PRIMARY KEY,
	    street VARCHAR(30) NOT NULL,
	    house INTEGER NOT NULL,
		building INTEGER NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initPositionsTable(db *sql.DB) error {
	env := "dbinit.initPositionsTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS positions(
	    id BIGSERIAL PRIMARY KEY,
	    position VARCHAR(30) NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initEventsTable(db *sql.DB) error {
	env := "dbinit.initEventsTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS events(
	    id BIGSERIAL PRIMARY KEY,
	    type_id BIGINT NOT NULL REFERENCES pricelist(id),
	    address_id BIGINT REFERENCES addresses(id),
		datetime DATETIME NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initPlayersTable(db *sql.DB) error {
	env := "dbinit.initPlayersTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS players(
	    id BIGSERIAL PRIMARY KEY,
	    name VARCHAR(30) NOT NULL,
		surname VARCHAR(30),
		tg_link VARCHAR(30) NOT NULL,
		is_sigma BOOL,
	    position_id BIGINT REFERENCES positions(id));
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_telegram ON players(tg_link);")
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_membership ON players(is_sigma);")
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initPlayerPresenceTable(db *sql.DB) error {
	env := "dbinit.initPlayerPresenceTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS player_presence(
	    id BIGSERIAL PRIMARY KEY,
	    event_id BIGINT NOT NULL REFERENCES events(id),
	    player_id BIGINT NOT NULL REFERENCES players(id));
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_presence ON player_presence(player_id);")
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func initPaymentsTable(db *sql.DB) error {
	env := "dbinit.initPaymentsTable"

	err := execStatement(db, `
	CREATE TABLE IF NOT EXISTS payments(
	    id BIGSERIAL PRIMARY KEY,
	    player_id BIGINT NOT NULL REFERENCES players(id),
		event_id BIGINT NOT NULL REFERENCES events(id),
		price INTEGER,
		payed BOOL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	err = execStatement(db, "CREATE INDEX IF NOT EXISTS idx_payments ON payments(player_id);")
	if err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func execStatement(db *sql.DB, query string) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error occured during preparation: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("error occured during execution: %w", err)
	}

	return nil
}
