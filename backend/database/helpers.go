// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/theparanoids/ashirt-server/backend/logging"

	"github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func (*Connection) FailIfTransaction(err error) {}

// Select executes the provided SelectBuilder query, and marshals the response into the provided slice.
// Note: this is for retriving multiple results, or "rows"
func (c *Connection) Select(modelSlice interface{}, sb squirrel.SelectBuilder) error {
	query, values, err := prepSquirrel(sb, nil)
	if err != nil {
		return err
	}
	return c.Sqlx.Select(modelSlice, query, values...)
}

// SelectRaw executes a raw SQL string
// Note: this is for retriving multiple results, or "rows"
func (c *Connection) SelectRaw(modelSlice interface{}, query string) error {
	return c.Sqlx.Select(modelSlice, query)
}

// Get executes the provided SelectBuilder query, and marshals the response into the provided structure.
// Note: this is for retriving a single value -- i.e. not a row, but a cell
func (c *Connection) Get(model interface{}, sb squirrel.SelectBuilder) error {
	query, values, err := prepSquirrel(sb, nil)
	if err != nil {
		return err
	}
	return c.Sqlx.Get(model, query, values...)
}

// Exec wraps sqlx.Exec, adding multi-value support and query logging
func (c *Connection) Exec(query string, values ...interface{}) error {
	newQuery, newValues, err := sqlx.In(query, values...)
	if err != nil {
		return err
	}
	logQuery(newQuery, newValues)
	_, err = c.Sqlx.Exec(newQuery, newValues...)
	return err
}

// Insert provides a generic way to insert a single row into the database. The function expects
// the table to insert into as well as a map of columnName : columnValue. Also assumes that
// created_at and updated_at columns exist (And are updated with the current time)
// Returns the inserted record id, or an error if the insert fails
func (c *Connection) Insert(tableName string, valueMap map[string]interface{}, onDuplicates ...interface{}) (int64, error) {
	ins := prepInsert(tableName, valueMap)

	ins, err := addDuplicatesClause(ins, onDuplicates...)
	if err != nil {
		return -1, err
	}

	result, err := c.execSquirrel(ins)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

// BatchInsert is similar to Insert, but instead is designed for multiple inserts. Note that only a single
// SQL query is run here.
//
// Parameters:
//
// tableName: the name of the target table
// count: the number of items needed to insert
// mapFn: A function that produces a single set of values for a new database row.
//
//	Note that this will be called <count> times
//
// Returns: an error if the insert fails
func (c *Connection) BatchInsert(tableName string, count int, mapFn func(int) map[string]interface{}, onDuplicates ...interface{}) error {
	if count == 0 {
		return nil
	}

	columns := []string{}
	query := squirrel.Insert(tableName)
	for idx := 0; idx < count; idx++ {
		valueMap := mapFn(idx)
		if len(valueMap) == 0 {
			continue
		}
		if len(columns) == 0 {
			for columnName := range valueMap {
				columns = append(columns, columnName)
			}
			query = query.Columns(columns...)
		}
		values := []interface{}{}
		for _, columnName := range columns {
			values = append(values, valueMap[columnName])
		}
		query = query.Values(values...)
	}

	query, err := addDuplicatesClause(query, onDuplicates...)
	if err != nil {
		return err
	}

	_, err = c.execSquirrel(query)
	return err
}

// Update executes the provided UpdateBuilder, in addition to the "updated_at" field, using the server time.
func (c *Connection) Update(ub squirrel.UpdateBuilder) error {
	ub = ub.Set("updated_at", time.Now())
	_, err := c.execSquirrel(ub)
	if err != nil {
		return fmt.Errorf("Unable to execute db update : %w", err)
	}
	return nil
}

// Delete removes records indicated by the given DeleteBuilder
func (c *Connection) Delete(dq squirrel.DeleteBuilder) error {
	_, err := c.execSquirrel(dq)
	return err
}

func (c *Connection) execSquirrel(sQuery squirrel.Sqlizer) (sql.Result, error) {
	query, values, err := prepSquirrel(sQuery, nil)
	if err != nil {
		return nil, err
	}
	return c.DB.Exec(query, values...)
}

func logQuery(query string, values []interface{}) {
	logging.SystemLog("msg", "executing query", "query", query, "values", fmt.Sprintf("%v", values))
}

// IsEmptyResultSetError returns true if the passed error is a database error resulting
// from querying a table expecting 1 row (with db.Get) but recieving 0
func IsEmptyResultSetError(err error) bool {
	return err == sql.ErrNoRows
}

// IsAlreadyExistsError returns true if the passed error is a database error resulting
// from attempting to insert a row with a unique/primary key that already exists
func IsAlreadyExistsError(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	return ok && mysqlErr.Number == 1062
}

func addDuplicatesClause(query squirrel.InsertBuilder, onDuplicates ...interface{}) (squirrel.InsertBuilder, error) {
	if len(onDuplicates) == 0 {
		return query, nil
	}
	stmt, ok := onDuplicates[0].(string)
	if !ok {
		return query, fmt.Errorf("onDuplicate[0] value must be a string")
	}
	if len(onDuplicates) > 1 {
		return query.Suffix(stmt, onDuplicates[1:]...), nil
	}
	return query.Suffix(stmt), nil
}
