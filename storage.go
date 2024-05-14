package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int64) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	user, dbname, password := os.Getenv("USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbname, password)
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
    encrypted_password VARCHAR(255) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (number)
)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {

	query := `INSERT INTO accounts (first_name, last_name, number, encrypted_password, balance, created_at)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM accounts WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetAccountByNumber(number int64) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts WHERE number = $1`, number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccounts(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", number)
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccounts(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	for rows.Next() {
		acc, err := scanIntoAccounts(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func scanIntoAccounts(rows *sql.Rows) (*Account, error) {
	acc := new(Account)
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt,
	)

	return acc, err
}
