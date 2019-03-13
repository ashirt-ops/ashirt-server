// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// Transactable is a wrapped sql.Tx (sql transaction). It provides two pieces of functionality.
// First: It allows execution of standard sql methods. These are all done in a transaction, on a
// single connection to the database -- the result being that either all of the queries succeed,
// or none of them do. Ultimately, this helps prevent the database from getting into an inconsistent
// state
// Second: this structure helps track errors that occur during the execution process. Any error that
// occurs during this processing is recorded, and can be retrieved by using the Error method. Internally,
// this is also used to determine whether we should commit the changes, or roll them back. Note that
// any error that occurs short-circuits any other call to be a no-op (except Rollback and Commit)
type Transactable struct {
	tx               *sqlx.Tx
	transactionError error
}

var errorNilTransaction = fmt.Errorf("transaction is nil")

// NewTransaction creates a new Transactable, which can then be used to execute queries in that
// transaction. Note that this is NOT the preferred method. Typically, the WithTx method should
// be used, as this will ensure that no db resources are lost.
func (c *Connection) NewTransaction(ctx context.Context) (*Transactable, error) {
	tx, err := c.Sqlx.BeginTxx(ctx, nil)
	rtn := Transactable{tx: tx, transactionError: err}
	return &rtn, err
}

// WithTx provides a scoped route to executing queries inside a transaction. Simply pass a function
// that uses the provided transaction to this method. At the end of the function, Commit will be
// called, which will either apply those actions to the database, or roll them back to the initial state
//
// Note 1: The provided function will not _necessarily_ be called. In situations where a transaction cannot
// be created, the function will be bypassed.
// Note 2: If the provided function needs to be exited early, you can call Rollback or Commit to ensure
// the desired outcome, and then return from the function as usual.
func (c *Connection) WithTx(ctx context.Context, fn func(tx *Transactable)) error {
	tx, _ := c.NewTransaction(ctx)
	if tx.Error() == nil {
		fn(tx)
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	return tx.Error()
}

// commit tries to execute the transaction, returning the error it received, if any
func (tx *Transactable) commit() error {
	if tx.tx == nil {
		return errorNilTransaction
	}
	return tx.tx.Commit()
}

// Rollback aborts the transaction, in effect rolling the database back to its initial state.
// Returns an error if the rollback encounters an error itself
func (tx *Transactable) Rollback() error {
	if tx.tx == nil {
		return errorNilTransaction
	}
	return tx.tx.Rollback()
}

// Commit attempts to apply the transaction to the database. If an error has been encountered, this
// method will instead attempt to rollback the transaction.
// An error is returned if any part of the transaction failed to execute, including the commit itself.
func (tx *Transactable) Commit() error {
	defer tx.Rollback() // A rollback attempt in case commit  fails.
	if tx.Error() != nil {
		return tx.Error()
	}
	tx.transactionError = tx.commit()
	return tx.Error()
}

// Error retrieves the recorded error, if any (nil otherwise)
func (tx *Transactable) Error() error {
	return tx.transactionError
}

// Select executes a SQL SELECT query, given the provided squirrel SelectBuilder and a reference to
// the database model. Results will be stored in the reference if any rows are returned.
func (tx *Transactable) Select(modelSlice interface{}, sb squirrel.SelectBuilder) error {
	return tx.sel(tx.tx.Select, modelSlice, sb)
}

// Get executes a SQL SELECT query given a reference to the database model and the squirrel SelectBuilder
// that returns a single row. If no rows are returned, or multiple rows are returned from the SQL
// query, an error is returned from this method.
func (tx *Transactable) Get(model interface{}, sb squirrel.SelectBuilder) error {
	return tx.sel(tx.tx.Get, model, sb)
}

// Insert executes a SQL INSERT query, given the tablename along with a columnName:value map.
// Returns the new row ID added (via LastInsertId), along with an error, if one was encountered
func (tx *Transactable) Insert(tableName string, valueMap map[string]interface{}) (int64, error) {
	if tx.Error() != nil {
		return -1, tx.Error()
	}

	ins := prepInsert(tableName, valueMap)
	res, err := tx.exec(ins)
	if err != nil {
		tx.transactionError = err
		return -1, err
	}
	id, err := res.LastInsertId()
	tx.transactionError = err
	return id, err
}

// BatchInsert executes a SQL INSERT query, for multiple value sets. This executes only a single query,
// but allows the caller to provide multiple db rows. The mapFn parameter should return back data for
// the i'th row to be added
// Returns an error if an error has been encountered.
func (tx *Transactable) BatchInsert(tableName string, count int, mapFn func(int) map[string]interface{}) error {
	if tx.Error() != nil {
		return tx.Error()
	}

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

	_, err := tx.exec(query)
	tx.transactionError = err
	return err
}

// Update executes a SQL UPDATE query, and sets the updated_at column to the server time
func (tx *Transactable) Update(ub squirrel.UpdateBuilder) error {
	ub = ub.Set("updated_at", time.Now())
	return tx.Exec(ub)
}

// Delete executes a SQL DELETE query
func (tx *Transactable) Delete(dq squirrel.DeleteBuilder) error {
	return tx.Exec(dq)
}

// Exec exectues any SQL query. Used primarily in situations where you need to interact with the
// dbms directly, or in rare situations where squirrel does not provide sufficient query modeling
// capability
func (tx *Transactable) Exec(s squirrel.Sqlizer) error {
	_, err := tx.exec(s)
	return err
}

// prepSquirrel does a ToSql call on the squirrel query, logs the query, then returns the result,
// or error if encountered
func prepSquirrel(s squirrel.Sqlizer, t *Transactable) (string, []interface{}, error) {
	if t != nil && t.tx == nil {
		return "", []interface{}{}, fmt.Errorf("transaction is nil")
	}

	query, values, err := s.ToSql()
	if err != nil {
		return "", []interface{}{}, err
	}
	newQuery, newValues, err := sqlx.In(query, values...)

	if err != nil {
		return "", []interface{}{}, err
	}

	logQuery(newQuery, newValues)

	return newQuery, newValues, err
}

// prepInsert "unzips" the columnName:value map into a set of columns and a set of values,
// then generates a squirrel InsertBuilder.
func prepInsert(tableName string, valueMap map[string]interface{}) squirrel.InsertBuilder {
	includeCreatedAt := true
	includeUpdatedAt := true

	columns := []string{}
	values := []interface{}{}

	for columnName, value := range valueMap {
		columns = append(columns, columnName)
		values = append(values, value)
		if columnName == "created_at" {
			includeCreatedAt = false
		} else if columnName == "updated_at" {
			includeUpdatedAt = false
		}
	}
	if includeCreatedAt {
		columns = append(columns, "created_at")
		values = append(values, time.Now())
	}
	if includeUpdatedAt {
		columns = append(columns, "updated_at")
		values = append(values, time.Now())
	}

	return squirrel.Insert(tableName).Columns(columns...).Values(values...)
}

// sel actually performs a select query, given the select method (either Sqlx.Select or Sqlx.Get)
func (tx *Transactable) sel(execFn func(interface{}, string, ...interface{}) error, model interface{}, sb squirrel.SelectBuilder) error {
	if tx.Error() != nil {
		return tx.Error()
	}
	query, values, err := prepSquirrel(sb, tx)
	if err != nil {
		tx.transactionError = err
		return err
	}
	err = execFn(model, query, values...)
	tx.transactionError = err
	return err
}

// exec provides the core logic for executing the Update/Delete/Exec public interface methods
func (tx *Transactable) exec(s squirrel.Sqlizer) (sql.Result, error) {
	if tx.Error() != nil {
		return nil, tx.Error()
	}

	query, values, err := prepSquirrel(s, tx)
	if err != nil {
		tx.transactionError = err
		return nil, err
	}
	res, err := tx.tx.Exec(query, values...)
	tx.transactionError = err
	return res, err
}
