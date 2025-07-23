package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type AccountStore interface { // interface since it defines the abstract behaviour of our store for Accounts
	CreateAccount(*Account) (*Account, error)
	DeleteAcount(int) error
	UpdateAccount(*Account) (*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct { // This will implmement the AccountStore interface. Go will implicitly know we implement it if it has all the required methods. Does not need an 'implements' or 'extends'
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
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
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL!")
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccount(*Account) (*Account, error) {
	account := NewAccount("New", "User")
	return account, nil
}

func (s *PostgresStore) UpdateAccount(*Account) (*Account, error) {
	account := NewAccount("Updated", "User")
	return account, nil
}

func (s *PostgresStore) DeleteAcount(int) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(int) (*Account, error) {
	account := NewAccount("ById", "User")
	return account, nil
}
