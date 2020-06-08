// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
)

// ListOperationsForAdmin is a specialized version of ListOperations where no operations are filtered
// For use in admin screens only
func ListOperationsForAdmin(ctx context.Context, db *database.Connection) ([]*dtos.OperationWithExportData, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var operations []struct {
		models.Operation
		NumUsers         int        `db:"num_users"`
		LastCompleteDate *time.Time `db:"last_complete_date"`
		ExportStatus     *int       `db:"most_recent_status"`
	}

	err := db.Select(&operations,
		sq.Select("id", "slug", "name", "ops.status", "count(user_id) AS num_users",
			"last_complete.max_date AS last_complete_date", "most_recent.status AS most_recent_status",
		).
			From("operations AS ops").
			LeftJoin("user_operation_permissions ON user_operation_permissions.operation_id = ops.id").
			JoinClause(
				"LEFT OUTER JOIN("+
					fmt.Sprintf("SELECT operation_id, MAX(updated_at) AS max_date FROM exports_queue WHERE status=%v GROUP BY operation_id", models.ExportStatusComplete)+
					") AS last_complete ON last_complete.operation_id = ops.id",
			).
			JoinClause(
				"LEFT OUTER JOIN("+
					"SELECT most_recent_event.operation_id, exports_queue.status FROM("+
					"SELECT operation_id, MAX(updated_at) AS max_date FROM exports_queue GROUP BY operation_id"+
					") AS most_recent_event "+
					"INNER JOIN exports_queue ON most_recent_event.operation_id = exports_queue.operation_id AND most_recent_event.max_date = exports_queue.updated_at"+
					") AS most_recent ON most_recent.operation_id = ops.id",
			).
			GroupBy("ops.id", "most_recent.status").
			OrderBy("ops.created_at DESC"))
	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	operationsDTO := []*dtos.OperationWithExportData{}
	for _, operation := range operations {
		operationsDTO = append(operationsDTO, &dtos.OperationWithExportData{
			Slug:         operation.Slug,
			Name:         operation.Name,
			Status:       operation.Status,
			NumUsers:     operation.NumUsers,
			ExportStatus: operation.ExportStatus,
			CompleteDate: operation.LastCompleteDate,

			// Temporary for screenshot client:
			ID: operation.ID,
		})
	}

	return operationsDTO, nil
}
