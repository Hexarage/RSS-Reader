package main

import (
	"database/sql"
	"log"
	"time"
)

type Storage interface {
	AddFeed(string) error
	DeleteFeed(int) error
	UpdateFeed(int, string) error // this really needed?
	GetFeedById(string) (string, error)
	GetAllFeeds() ([]string, error)
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

func (d *PostgresStore) Init() error {
	return d.CreateLinkTable()
}

func (d *PostgresStore) CreateLinkTable() error {
	query := `create table if not exists links (
		id serial primary key,
		link varcha(100),
		number serial,
		created_at timestamp
	)` // no clue if this is the correct way to do this
	// TODO: Fix this create table querry

	_, err := d.db.Exec(query)
	return err
}

func (d *PostgresStore) AddFeed(link string) error {

	query := `insert into links
	(link, number, created_at)
	values ($1, $2, $3)
	`

	response, err := d.db.Query(query, link, 0 /*no clue if this is even needed*/, time.Now().UTC()) // technically, this may not be considered correct, time of creation should be request received time
	if err != nil {
		return err
	}

	log.Printf("%+v\n", response)

	return nil
}

func (d *PostgresStore) DeleteFeed(id int) error {
	return nil
}

func (d *PostgresStore) UpdateFeed(id int, link string) error {
	return nil
}

func (d *PostgresStore) GetFeedById(string) (string, error) {
	return "", nil
}

func (d *PostgresStore) GetAllFeeds() ([]string, error) {

	rows, err := d.db.Query(`select * from links`)
	if err != nil {
		return nil, err
	}

	var links []string
	for rows.Next() {
		var link string
		var number int
		var createdAt time.Time
		if err := rows.Scan(&link, &number, &createdAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}
