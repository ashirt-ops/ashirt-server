// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package database

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

// Connection contains the infrastructure needed to manage the database connection
type Connection struct {
	DB                *sql.DB
	Sqlx              *sqlx.DB
	migrationsDirPath string
}

// NewConnection establishes a new connection to the databse server
func NewConnection(dsn string, migrationsDirPath string) (*Connection, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	return newConnection(config, migrationsDirPath)
}

func newConnection(config *mysql.Config, migrationsDirPath string) (*Connection, error) {
	config.ParseTime = true

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	sqlxDB := sqlx.NewDb(db, "mysql")

	return &Connection{db, sqlxDB, migrationsDirPath}, nil
}

// CheckSchema checks the database schema against the migrations
// and returns an error if they don't match
func (c *Connection) CheckSchema() error {
	m, err := c.migrationSource()
	if err != nil {
		return err
	}
	allMigrations, err := m.FindMigrations()
	if err != nil {
		return err
	}
	ranMigrations, err := migrate.GetMigrationRecords(c.DB, "mysql")
	if err != nil {
		return err
	}

	if len(ranMigrations) != len(allMigrations) {
		var suggestedFix string
		if len(ranMigrations) > len(allMigrations) {
			suggestedFix = "Check that you have deployed the latest version"
		} else {
			suggestedFix = "Run migrate up"
		}
		return fmt.Errorf(
			"Database has %d applied migrations but codebase expects %d migrations. %s",
			len(ranMigrations), len(allMigrations), suggestedFix,
		)
	}

	for i := range ranMigrations {
		if ranMigrations[i].Id != allMigrations[i].Id {
			return fmt.Errorf(
				"Ran migration id %s and codebase migration id %s do not match!",
				ranMigrations[i].Id, allMigrations[i].Id,
			)
		}
	}
	return nil
}

// MigrateUp performs all migrations
func (c *Connection) MigrateUp() error {
	m, err := c.migrationSource()
	if err != nil {
		return err
	}
	_, err = migrate.Exec(c.DB, "mysql", m, migrate.Up)
	return err
}

func (c *Connection) migrationSource() (migrate.MigrationSource, error) {
	migrationSource := &migrate.FileMigrationSource{Dir: c.migrationsDirPath}
	return migrationSource, nil
}
