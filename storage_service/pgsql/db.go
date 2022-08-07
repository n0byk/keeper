package pgsql

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Database struct {
	db  *pgx.Conn
	ctx context.Context
}

// NewDatabase creates and initializes a new Database
func NewDatabase() *Database {
	db := &Database{ctx: context.Background()}
	db.init()

	return db
}

// Init connects to a PostgreSQL instance and creates the
// tables this service relies on if they don't already exist
func (d *Database) init() {
	// read the PostgreSQL connection info from the environment
	config, err := pgx.ParseConfig("postgres://keeper:keeper@localhost:5432/keeper?sslmode=disable")
	if err != nil {
		panic(err)
	}

	// connect to database using configuration created from environment variables
	if d.db, err = pgx.ConnectConfig(d.ctx, config); err != nil {
		panic(err)
	}

	// create tables
	d.createTables()
}

// Close closes connections to the database
func (d *Database) Close() error {
	return d.db.Close(d.ctx)
}

// createTables creates the tables that this service relies on
func (d *Database) createTables() {
	query := `
 	create extension if not exists "uuid-ossp";
	CREATE TABLE IF NOT EXISTS public."user_catalog" (
		"user_id" uuid NOT NULL default uuid_generate_v4(),
		"user_login" varchar NOT NULL,
		"user_password" varchar NOT NULL, 
		"user_token" varchar, 
		"add_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"update_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"delete_time" timestamp, 
		CONSTRAINT "user_catalog_pk" PRIMARY KEY ("user_id")
	);

	CREATE TABLE IF NOT EXISTS public."user_data" (
		"data_id" uuid NOT NULL default uuid_generate_v4(),
		"user_id" varchar NOT NULL, 
		"row_data" varchar NOT NULL, 
		"add_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"update_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"delete_time" timestamp, 
		CONSTRAINT "user_data_pk" PRIMARY KEY ("data_id") 
	);

	`

	if _, err := d.db.Exec(d.ctx, query); err != nil {
		panic(err)
	}
}

// Add new User
func (d *Database) AddUser(login, password string) error {
	query := `
	INSERT INTO user_catalog(user_login, user_password) VALUES($1, $2)`

	_, err := d.db.Exec(d.ctx, query, login, password)
	return err
}

// Set user Token to user
func (d *Database) SetToken(user_id, user_token string) error {
	query := `UPDATE user_catalog set user_token = $2 where user_id = $1)`

	_, err := d.db.Exec(d.ctx, query, user_id, user_token)
	return err
}

// Validate user credentials
func (d *Database) FindUser(login, password string) (bool, error) {
	query := `SELECT true FROM user_catalog WHERE user_login = $1, user_password = $2 and delete_time is null`

	var val bool
	return val, d.db.QueryRow(d.ctx, query, login, password).Scan(&val)
}
