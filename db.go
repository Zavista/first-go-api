package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type AccountStore interface { // interface since it defines the abstract behaviour of our store for Accounts
	CreateAccount(*CreateAccountRequest) (*Account, error)
	DeleteAccount(int) error
	UpdateAccount(int, *UpdateAccountRequest) (*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountBalanceByID(int) (int64, error)
}

type PostgresStore struct { // This will implmement the AccountStore interface. Go will implicitly know we implement it if it has all the required methods. Does not need an 'implements' or 'extends'
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) { // Constructor Function
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, name,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL!")
	return &PostgresStore{
		db: db,
	}, nil
}

// Setup initializes the accounts table and triggers
func (s *PostgresStore) Setup() error {
	if err := s.createAccountTable(); err != nil {
		return err
	}
	if err := s.createUpdatedAtTrigger(); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number SERIAL,
		balance BIGINT DEFAULT 0,
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createUpdatedAtTrigger() error {
	fn := `
	CREATE OR REPLACE FUNCTION set_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = now();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;
	`
	tr := `
	CREATE TRIGGER trigger_set_updated_at
	BEFORE UPDATE ON accounts
	FOR EACH ROW
	EXECUTE FUNCTION set_updated_at();
	`

	if _, err := s.db.Exec(fn); err != nil {
		return err
	}
	if _, err := s.db.Exec(tr); err != nil {
		// ignore "already exists" errors silently
		if err.Error() != `pq: trigger "trigger_set_updated_at" for relation "accounts" already exists` {
			return err
		}
	}
	return nil
}

func (s *PostgresStore) CreateAccount(req *CreateAccountRequest) (*Account, error) {
	query := `
		INSERT INTO accounts (first_name, last_name)
		VALUES ($1, $2)
		RETURNING id, first_name, last_name, number, balance, created_at, updated_at;
	`

	row := s.db.QueryRow(query, req.FirstName, req.LastName)

	var created Account
	err := row.Scan(
		&created.ID,
		&created.FirstName,
		&created.LastName,
		&created.Number,
		&created.Balance,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (s *PostgresStore) UpdateAccount(id int, req *UpdateAccountRequest) (*Account, error) {
	query := `
		UPDATE accounts
		SET first_name = $1, last_name = $2, balance = $3
		WHERE id = $4
		RETURNING id, first_name, last_name, number, balance, created_at, updated_at;
	`

	row := s.db.QueryRow(query, req.FirstName, req.LastName, req.Balance, id)

	var updated Account
	err := row.Scan(
		&updated.ID,
		&updated.FirstName,
		&updated.LastName,
		&updated.Number,
		&updated.Balance,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = $1;`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `
		SELECT id, first_name, last_name, number, balance, created_at, updated_at
		FROM accounts
		WHERE id = $1;
	`

	row := s.db.QueryRow(query, id)

	var acc Account
	err := row.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Balance,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no account found with id %d", id)
		}
		return nil, err
	}

	return &acc, nil
}

func (s *PostgresStore) GetAccountBalanceByID(id int) (int64, error) {
	query := `SELECT balance FROM accounts WHERE id = $1;`

	var balance int64
	err := s.db.QueryRow(query, id).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no account found with id %d", id)
		}
		return 0, err
	}

	return balance, nil
}
