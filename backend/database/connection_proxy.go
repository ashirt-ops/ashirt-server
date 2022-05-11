// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package database

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

// ConnectionProxy provides an interface into the database, using either an underlying connection,
// or a transaction. This is compatible with both Connection and Transactable types
type ConnectionProxy interface {
	Select(modelSlice interface{}, builder sq.SelectBuilder) error
	Get(model interface{}, builder sq.SelectBuilder) error
	// Exec(query string, values ...interface{}) error // not consistent between Transactable and Connection
	Insert(table string, valueMap map[string]interface{}, onDuplicates ...interface{}) (int64, error)
	BatchInsert(tableName string, count int, mapFn func(int) map[string]interface{}, onDuplicates ...interface{}) error
	Update(builder sq.UpdateBuilder) error
	Delete(builder sq.DeleteBuilder) error
	WithTx(ctx context.Context, fn func(tx *Transactable)) error
}

// _verifyConnectionProxyInterface is a "canary" function that ensures that the expected concrete
// types to the ConnectionProxy interface properly implement the interface.
//lint:ignore U1000 This is just to verify the interface -- it should never be called directly anyway
func _verifyConnectionProxyInterface() {
	var conn *Connection
	var tx *Transactable
	check := func(c ConnectionProxy) {}
	check(conn)
	check(tx)
}
