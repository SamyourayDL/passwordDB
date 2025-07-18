package postgres

import (
	"database/sql"
	"password-db/internals/config"

	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	db  *sql.DB
	key string
}

type PassResponse struct {
	Service  string
	Password string
}

func New(cfg *config.DB) (*Storage, error) {
	const fn = "internals.storage.postgres.New"

	databaseUrl := "postgres://" + cfg.User + ":" + cfg.Password + "@" + cfg.Address + "/" + cfg.Name + "?sslmode=disable"
	// fmt.Println("url: ", databaseUrl)
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{
		db:  db,
		key: cfg.Key,
	}, nil
}

func (s *Storage) GetPass(user, service string) ([]PassResponse, error) {
	const fn = "internals.storage.postgres.GetPass"

	var queryConstraint string
	if service != "" {
		queryConstraint = " AND service_name == ?"
	}

	stmt, err := s.db.Prepare(`
		SELECT service_name, secret_enc
		FROM passwords
		JOIN users
		ON passwords.user_id = users.id
		WHERE name == ?
	` + queryConstraint)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

	if service != "" {
		stmt.Exec()
	}
}
