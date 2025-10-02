package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/macadamiaboy/SigmaPay/internal/config"
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

	if err = InitDatabase(db); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
