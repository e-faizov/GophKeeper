package stores

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog/log"
)

func tableExist(ctx context.Context, db *sql.DB, tb string) bool {
	sql := `select exists (
	   select from information_schema.tables
	   where  table_schema = 'public'
	   and    table_name   = $1
	   )
`

	var exists bool
	row := db.QueryRowContext(ctx, sql, tb)
	err := row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func createTable(ctx context.Context, db *sql.DB, sqls ...string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error open Tx: %w", err)
	}
	rollback := func(err error) error {
		errRoll := tx.Rollback()
		if errRoll != nil {
			err = multierror.Append(err, fmt.Errorf("error on rollback %w", errRoll))
		}
		return err
	}
	for _, s := range sqls {
		_, err = tx.ExecContext(ctx, s)
		if err != nil {
			return rollback(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func initTables(ctx context.Context, db *sql.DB) error {
	var err error
	exist := tableExist(ctx, db, "users")
	if !exist {
		err = createUsersTable(ctx, db)
		if err != nil {
			return fmt.Errorf("error create users: %v", err)
		}
		log.Info().Msg("users table created")
	}

	exist = tableExist(ctx, db, "secrets")
	if !exist {
		err = createSecretsTable(ctx, db)
		if err != nil {
			return fmt.Errorf("error create secrets: %v", err)
		}
		log.Info().Msg("secrets table created")
	}

	return nil
}

func createUsersTable(ctx context.Context, db *sql.DB) error {
	err := createTable(ctx, db,
		`create table users
(
	id serial,
	uuid text,
	login text,
	hash text
)`,
		`create unique index users_uuid_uindex
	on users (uuid)`,
		`create unique index users_login_uindex
	on users (login)`,
		`alter table users
	add constraint users_pk
		primary key (login)`)
	return err
}

func createSecretsTable(ctx context.Context, db *sql.DB) error {
	err := createTable(ctx, db,
		`create table secrets
(
	id serial,
	uid text,
	user_uid text,
	data1 text,
	data2 text,
	data3 text,
	version int,
	type int
)`,
		`create unique index secrets_data1_uindex
	on secrets (data1)`,
		`create unique index secrets_id_uindex
	on secrets (id)`,
		`create unique index secrets_uid_uindex
	on secrets (uid)`,
		`alter table secrets
	add constraint secrets_pk
		primary key (uid)`,
	)
	return err
}
