package app

import (
	"context"

	"github.com/google/uuid"
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

// NewDatabase creates and initializes a new Database
func Connection() *Database {
	db := &Database{ctx: context.Background()}
	db.init()

	return db
}

// Init connects to a PostgreSQL instance and creates the
// tables this service relies on if they don't already exist
func (d *Database) init() {
	// read the PostgreSQL connection info from the environment
	config, err := pgx.ParseConfig("postgres://keeper:keeper@keeper_postgres:5432/keeper?sslmode=disable")
	if err != nil {
		panic(err)
	}

	// connect to database using configuration created from environment variables
	if d.db, err = pgx.ConnectConfig(d.ctx, config); err != nil {
		panic(err)
	}

}

// Close closes connections to the database
func (d *Database) Close() error {
	return d.db.Close(d.ctx)
}

// createTables creates the tables that this service relies on
func (d *Database) CreateTables() {
	query := `
 	create extension if not exists "uuid-ossp";
	CREATE TABLE IF NOT EXISTS public."user_catalog" (
		"user_id" uuid NOT NULL default uuid_generate_v4(),
		"user_login" varchar NOT NULL,
		"user_password" varchar NOT NULL, 
		"user_token" varchar, 
		"public_token"  uuid NOT NULL default uuid_generate_v4()  UNIQUE, 
		"add_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"update_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"delete_time" timestamp, 
		CONSTRAINT "user_catalog_pk" PRIMARY KEY ("user_id")
	);

	CREATE TABLE IF NOT EXISTS public."user_data" (
		"data_id" uuid NOT NULL default uuid_generate_v4(),
		"user_id" varchar NOT NULL, 
		"row_data" varchar NOT NULL, 
		"lobby_id" uuid NOT NULL REFERENCES lobby_catalog (lobby_id), 
		"add_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"update_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"delete_time" timestamp, 
		CONSTRAINT "user_data_pk" PRIMARY KEY ("data_id") 
	);

	CREATE TABLE IF NOT EXISTS public."lobby_catalog" (
		"lobby_id" uuid NOT NULL default uuid_generate_v4(),
		"lobby_name" varchar NOT NULL, 
		"owner_user_id" uuid NOT null REFERENCES user_catalog (user_id), 
		"add_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"update_time" timestamp NOT NULL DEFAULT (now() at time zone 'UTC'),
		"delete_time" timestamp, 
		CONSTRAINT "lobby_catalog_pk" PRIMARY KEY ("lobby_id")
	);
	
	CREATE TABLE IF NOT EXISTS public."lobby_user_set" (
		"set_id" uuid NOT NULL default uuid_generate_v4(),
		"lobby_id" uuid NOT NULL REFERENCES lobby_catalog (lobby_id), 
		"public_token" uuid NOT NULL REFERENCES user_catalog (public_token), 
		CONSTRAINT "lobby_user_set_pk" PRIMARY KEY ("set_id")
	);

	create unique index IF NOT EXISTS user_login_unique_idx on user_catalog (user_login);
	create unique index IF NOT EXISTS lobby_name_unique_idx on lobby_catalog (lobby_name);
	`

	if _, err := d.db.Exec(d.ctx, query); err != nil {
		panic(err)
	}
}

// Add new User
func (d *Database) AddUser(data TRegistrationRequest) error {
	query := `INSERT INTO user_catalog(user_login, user_password) VALUES($1, $2);`

	_, err := d.db.Exec(d.ctx, query, data.Login, HashPassword(data.Password))
	return err
}

func (d *Database) AddLobby(data TCreateLobbyRequest) error {

	lobbyId := uuid.Must(uuid.New(), nil).String()

	query := `
	INSERT INTO lobby_catalog (lobby_name, owner_user_id, lobby_id) SELECT $1, uc.user_id, $3 FROM user_catalog uc WHERE uc.user_token = $2; 
`

	_, err := d.db.Exec(d.ctx, query, data.LobbyName, data.Token, lobbyId)
	if err != nil {
		return err
	}

	query2 := ` 
	INSERT INTO lobby_user_set (lobby_id, public_token) SELECT $1, uc2.public_token FROM user_catalog uc2 WHERE uc2.user_token = $2;
`
	_, err = d.db.Exec(d.ctx, query2, lobbyId, data.Token)

	return err
}

// Set user Token to user
func (d *Database) SetToken(user_login, user_token string) error {
	query := `UPDATE user_catalog set user_token = $2 where user_login = $1;`

	_, err := d.db.Exec(d.ctx, query, user_login, user_token)
	return err
}

// Validate user credentials
func (d *Database) FindUser(login, password string) (bool, error) {
	query := `SELECT true FROM user_catalog WHERE user_login = $1 and user_password = $2 and delete_time is null;`

	var val bool
	return val, d.db.QueryRow(d.ctx, query, login, HashPassword(password)).Scan(&val)
}

func (d *Database) ValidateToken(token string) (bool, error) {
	query := `SELECT true FROM user_catalog WHERE user_token = $1 and delete_time is null;`

	var val bool
	return val, d.db.QueryRow(d.ctx, query, token).Scan(&val)
}

func (d *Database) SetUserToLobby(data TAddToLobbyRequest) error {

	query := `INSERT INTO lobby_user_set (lobby_id, public_token)(SELECT  lobby_id, $2 FROM lobby_catalog WHERE lobby_name = $1)`

	_, err := d.db.Exec(d.ctx, query, data.LobbyName, data.PublicToken)
	return err
}

func (d *Database) findPublicToken(login string) (string, error) {
	query := `SELECT public_token FROM user_catalog WHERE user_login = $1 and delete_time is null;`

	var publicToken string
	return publicToken, d.db.QueryRow(d.ctx, query, login).Scan(&publicToken)
}

func (d *Database) findLobbyId(name string) (string, error) {
	query := `SELECT lobby_id FROM public.lobby_catalog where lobby_name = $1`

	var lobbyId string
	return lobbyId, d.db.QueryRow(d.ctx, query, name).Scan(&lobbyId)
}

func (d *Database) ValidateLobbyPermission(data TValidatePermissions) (string, error) {

	query := `select lc.lobby_name from lobby_catalog lc 
	left join user_catalog uc on uc.user_token = $1
	left join lobby_user_set lus on uc.public_token = lus.public_token 
	where lc.lobby_name = $2`

	var lobby string
	return lobby, d.db.QueryRow(d.ctx, query, data.Token, data.LobbyName).Scan(&lobby)
}

func (d *Database) AddMessage(data TSetDataRequest) (string, error) {

	query := `INSERT INTO user_data (user_id, row_data, lobby_id)(SELECT  user_id, $2, $1 FROM user_catalog  WHERE user_token  = $3)`

	_, err := d.db.Exec(d.ctx, query, data.LobbyId, data.RowData, data.Token)
	return data.RowData, err
}
