package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=qsmsoft dbname=gobank password=cs1993adm sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()

}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    number BIGINT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (number)
)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {

	query := `INSERT INTO accounts (first_name, last_name, number, balance, created_at)
    VALUES ($1, $2, $3, $4, $5)`

	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", resp)

	return nil
}
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) GetAccountID(id int) (*Account, error) {
	row, err := s.db.Query(`SELECT * FROM accounts WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	acc := new(Account)

	if err := row.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Balance,
		&acc.CreatedAt,
	); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		acc := new(Account)
		if err := rows.Scan(
			&acc.ID,
			&acc.FirstName,
			&acc.LastName,
			&acc.Number,
			&acc.Balance,
			&acc.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}
