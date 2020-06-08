package services

import (
	"context"

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
	err = db.Exec("INSERT INTO exports_queue(operation_id, status, user_id)" +
	" SELECT ?, ?, ? FROM DUAL WHERE NOT EXISTS(" +
	" SELECT id FROM exports_queue WHERE operation_id=? AND status=? LIMIT 1);",
	op.ID, models.ExportStatusPending, middleware.UserID(ctx), op.ID, models.ExportStatusPending)

	if err != nil {
		return errRtn(backend.DatabaseErr(err))
	}
	return ExportOperationOutput{Queued: queued}
}
