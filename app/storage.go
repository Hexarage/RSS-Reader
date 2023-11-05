package main

import "database/sql"

type Storage interface {
	AddFeed(string) error
	DeleteFeed(int) error
	UpdateFeed(string) error // this really needed?
	GetFeedById(int) (string, error)
}

type PostgresStore struct {
	db sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connectionString := "user=postgresuser dbname=postgres password=postgrespassword sslmode=disable" //hardcoding the info; obviously sslmode would be verify-full in proper live environment
	dbl, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := dbl.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: *dbl,
	}, nil
}

func (d *PostgresStore) AddFeed(string) error {
	return nil
}

func (d *PostgresStore) DeleteFeed(int) error {
	return nil
}

func (d *PostgresStore) UpdateFeed(string) error { // this really needed?
	return nil
}

func (d *PostgresStore) GetFeedById(int) (string, error) {
	return "", nil
}
