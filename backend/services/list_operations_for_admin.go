// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
)

// ListOperationsForAdmin is a specialized version of ListOperations where no operations are filtered
// For use in admin screens only
func ListOperationsForAdmin(ctx context.Context, db *database.Connection) ([]*dtos.Operation, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}
	return listAllOperations(db)
}
