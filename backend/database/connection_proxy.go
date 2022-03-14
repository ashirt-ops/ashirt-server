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
	// Exec(query string, values ...interface{}) error // not consistent between Transactable and
	Insert(table string, valueMap map[string]interface{}) (int64, error)
	BatchInsert(tableName string, count int, mapFn func(int) map[string]interface{}) error
	Update(builder sq.UpdateBuilder) error
	Delete(builder sq.DeleteBuilder) error
	WithTx(ctx context.Context, fn func(tx *Transactable)) error
}
