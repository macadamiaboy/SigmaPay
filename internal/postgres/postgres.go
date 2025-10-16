package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/macadamiaboy/SigmaPay/internal/config"
)

type DataBase struct {
	Connection *sql.DB
}

func New() error {
	const env = "postgres.New"

	db, err := PrepareDB()
	if err != nil {
		log.Fatalf("%s: failed to prepare the db: %v", env, err)
		return fmt.Errorf("%s: %w", env, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	if err = InitDatabase(db.Connection); err != nil {
		return fmt.Errorf("%s: %w", env, err)
	}

	return nil
}

func PrepareDB() (*DataBase, error) {
	const env = "postgres.PrepareDB"

	pgConfig := config.LoadDBConfigData()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable",
		pgConfig.Database.Username,
		pgConfig.Database.Password,
		pgConfig.Database.Host,
		pgConfig.Database.Port,
		pgConfig.Database.DBName,
	)

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", env, err)
	}

	return &DataBase{Connection: conn}, nil
}

func (db *DataBase) Close() error {
	return db.Connection.Close()
}
