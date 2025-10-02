package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/macadamiaboy/SigmaPay/config"
)

type Storage *sql.DB

func New() error {
	const op = "postgres.New"

	pgConfig := config.LoadDBConfigData()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		pgConfig.Database.Username,
		pgConfig.Database.Password,
		pgConfig.Database.Host,
		pgConfig.Database.Port,
		pgConfig.Database.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS pricelist(
	    id INTEGER PRIMARY KEY,
	    type VARCHAR(20) NOT NULL,
	    price INTEGER NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE INDEX IF NOT EXISTS idx_type ON pricelist(type);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS addresses(
	    id INTEGER PRIMARY KEY,
	    street VARCHAR(30) NOT NULL,
	    house INTEGER NOT NULL,
		building INTEGER NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS positions(
	    id INTEGER PRIMARY KEY,
	    position VARCHAR(30) NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS events(
	    id INTEGER PRIMARY KEY,
	    type_id INTEGER NOT NULL REFERENCES pricelist(id),
	    address_id INTEGER NOT NULL REFERENCES addresses(id),
		time TIME NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS players(
	    id INTEGER PRIMARY KEY,
	    name VARCHAR(30) NOT NULL,
		surname VARCHAR(30),
		tg_link VARCHAR(30) NOT NULL,
		is_sigma BOOL,
	    position_id INTEGER REFERENCES positions(id));
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE INDEX IF NOT EXISTS idx_type ON players(tg_link);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE INDEX IF NOT EXISTS idx_type ON players(is_sigma);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS player_presence(
	    id INTEGER PRIMARY KEY,
	    event_id INTEGER NOT NULL REFERENCES events(id),
	    player_id INTEGER NOT NULL REFERENCES players(id),
		price INTEGER NOT NULL);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	/*
		if err = db.Ping(); err != nil {
			fmt.Println(charmap.Windows1251.NewDecoder().String(err.Error()))
			return fmt.Errorf("%s: %w", op, err)
		}
	*/

	return nil
}
