package services

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"
)

// ExportOperationOutput indicates the action (or inaction) performed as a result of the request.
// 1. If Err is not nil, then the request did not succeed -- do not act on the rest of the response (failure)
// 2. If Queued is true, it means that a new item has been added to the queue (Success)
// 3. If Queued is false, it means that this is already in-queue, and no action has been taken (also a success)
type ExportOperationOutput struct {
	Err    error
	Queued bool
}

// CreateOperationExport marks an operation for (eventual) export. The actual export occurs in processors
func CreateOperationExport(ctx context.Context, db *database.Connection, operationSlug string) ExportOperationOutput {
	errRtn := func(er error) ExportOperationOutput { return ExportOperationOutput{Err: er} } // re-package for simplier error responses

	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return errRtn(backend.UnauthorizedWriteErr(err))
	}

	op, err := lookupOperation(db, operationSlug)
	if err != nil {
		return errRtn(backend.DatabaseErr(err))
	}

	queued := false
	err = db.WithTx(ctx, func(tx *database.Transactable) {
		// TODO: there is some potential here for multiple, duplicate rows to be inserted (during concurrent requests)
		// this probably cannot be resolved via mysql's replace or "on duplicate key update". However, stored procedures might work
		// however, this should be unlikely given the expected low-usage of this feature (i.e. exports should be few and far between)
		var count int64
		tx.Get(&count, sq.Select("count(id)").From("exports_queue").
			Where(sq.Eq{"exports_queue.operation_id": op.ID}).
			Where(sq.Eq{"exports_queue.status": models.ExportStatusPending}),
		)
		if count == 0 {
			tx.Insert("exports_queue", map[string]interface{}{
				"operation_id": op.ID,
				"user_id":      middleware.UserID(ctx),
				"status":       models.ExportStatusPending,
			})
			queued = true
		}
	})

	if err != nil {
		return errRtn(backend.DatabaseErr(err))
	}
	return ExportOperationOutput{Queued: queued}
}
