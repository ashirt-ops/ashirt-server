// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

// NewTestConnection creates a new empty database, and applies for integration/service tests
func NewTestConnection(t *testing.T, dbName string) *Connection {
	return NewTestConnectionFromNonStandardMigrationPath(t, dbName, "../migrations")
}

// NewTestConnectionFromNonStandardMigrationPath is identical to NewTestConnection, but allows the user to specify where the migrations dir is
func NewTestConnectionFromNonStandardMigrationPath(t *testing.T, dbName, migrationsDirPath string) *Connection {
	config := &mysql.Config{
		User:   "root",
		Passwd: "dev-root-password",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: dbName,
	}

	setupDBConfig := &mysql.Config{
		User:            config.User,
		Passwd:          config.Passwd,
		Net:             config.Net,
		Addr:            config.Addr,
		MultiStatements: true,
	}

	setupDB, err := sql.Open("mysql", setupDBConfig.FormatDSN())
	require.NoError(t, err)

	_, err = setupDB.Exec("CREATE DATABASE IF NOT EXISTS `" + config.DBName + "`")
	require.NoError(t, err)

	// Reset database back to schema.sql
	schemaSQL, err := os.ReadFile(migrationsDirPath + "/../schema.sql")
	require.NoError(t, err)

	_, err = setupDB.Exec("USE `" + config.DBName + "`;\n" + string(schemaSQL))
	require.NoError(t, err)

	db, err := newConnection(config, migrationsDirPath)
	require.NoError(t, err)

	return db
}
