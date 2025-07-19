package postgres

import (
	"database/sql"
	"password-db/internals/config"
	"password-db/internals/lib/crypto"
	"time"

	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	db  *sql.DB
	key []byte
}

type PassResponse struct {
	Service  string
	Password string
}

const EmptyCategory = "-"

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

	queryConstraint := ""
	if service != "" {
		queryConstraint = " AND service_name = $2"
	}
	q := `SELECT service_name, secret_enc
		FROM passwords
		JOIN users
		ON passwords.user_id=users.id
		WHERE users.name = $1` + queryConstraint

	var (
		rows *sql.Rows
		err  error
	)

	if service != "" {
		rows, err = s.db.Query(q, user, service)
	} else {
		rows, err = s.db.Query(q, user)
	}

	fmt.Println("Query was executed with:", err)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var result []PassResponse
	for rows.Next() {
		var (
			service    string
			cipherpass []byte
		)

		if err := rows.Scan(&service, &cipherpass); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		pass, err := crypto.Decrypt(cipherpass, s.key)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		result = append(result, PassResponse{
			Service:  service,
			Password: string(pass),
		})
	}

	return result, nil
}

func (s *Storage) AddUser(user string) error {
	const fn = "internals.storage.postgres.AddUser"

	q := `INSERT INTO users (name) VALUES ($1)`

	if _, err := s.db.Exec(q, user); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) AddPassword(user, password, service_name, category string) error {
	const fn = "internals.storage.postgres.AddPassword"

	qUserId := `SELECT id FROM users WHERE name=$1`

	res := s.db.QueryRow(qUserId, user)
	var userId int
	if err := res.Scan(&userId); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if category == "" {
		category = EmptyCategory
	}

	q := `INSERT INTO passwords (user_id, service_name, secret_enc, category, created_at)
	VALUES
		($1, $2, $3, $4, $5)`

	encPassword, err := crypto.Encrypt([]byte(password), s.key)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if _, err := s.db.Exec(q, userId, service_name, encPassword, category, time.Now()); err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) Delete(user, service string) {}
