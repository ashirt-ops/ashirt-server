// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"testing"

	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"
)

func internalTestDBSetup(t *testing.T) *database.Connection {
	return database.NewTestConnection(t, "service-test-db")
}

func mkTestingContext(userid int64, p policy.Policy) context.Context {
	ctx := context.Background()
	return middleware.InjectIntoContext(ctx, middleware.InjectIntoContextInput{
		UserID:       userid,
		IsSuperAdmin: true,
		UserPolicy:   p,
	})
}

type mockOperation struct {
	User     dtos.User
	UserID   int64
	Op       *dtos.Operation
	ID       int64
	Findings []models.Finding
	Evidence []models.Evidence
}

// setupBasicTestOperation creates two opeartions in the database, with evidence and findings (unassociated)
// the first operation has a small collection of evidence and findings and is intended to be tested against
// the second operation has a single piece of evidence and finding, to do tests for "such and such does not belong to this operation" branches
// All wiring still needs to be done by the user
func setupBasicTestOperation(t *testing.T, db *database.Connection) (mockOperation, mockOperation) {

	commonFindingCategories := []string{
		"Product",
		"Network",
		"Enterprise",
		"Vendor",
		"Behavioral",
		"Detection Gap",
	}
	db.BatchInsert("finding_categories", len(commonFindingCategories), func(i int) map[string]interface{} {
		return map[string]interface{}{
			"category": commonFindingCategories[i],
		}
	})

	goodOp := mockOperation{
		User:     dtos.User{FirstName: "fn", LastName: "ln", Slug: "sn"},
		Findings: make([]models.Finding, 0),
		Evidence: make([]models.Evidence, 0),
	}
	badOp := mockOperation{
		User:     dtos.User{FirstName: "fn", LastName: "ln", Slug: "sn"},
		Findings: make([]models.Finding, 0),
		Evidence: make([]models.Evidence, 0),
	}
	userID, err := db.Insert("users", map[string]interface{}{
		"slug":       goodOp.User.Slug,
		"first_name": goodOp.User.FirstName,
		"last_name":  goodOp.User.LastName,
		"email":      "",
	})
	require.NoError(t, err)
	goodOp.UserID = userID
	badOp.UserID = userID
	ctx := mkTestingContext(userID, &policy.FullAccess{})

	makeOp := func(slug, name string) (*dtos.Operation, int64) {
		op, err := CreateOperation(ctx, db, CreateOperationInput{Slug: slug, OwnerID: userID, Name: name})
		require.NoError(t, err)

		var opID int64
		err = db.Get(&opID, sq.Select("id").From("operations").Where(sq.Eq{"slug": op.Slug}))
		require.NoError(t, err)

		return op, opID
	}


	goodOp.Op, goodOp.ID = makeOp("goodOp", "Good Operation")
	badOp.Op, badOp.ID = makeOp("badOp", "Bad Operation")

	cs, _ := contentstore.NewMemStore()
	makeEvidence := func(op *mockOperation, id int64, desc string) {
		input := CreateEvidenceInput{OperationSlug: op.Op.Slug, Description: desc, ContentType: "other"}
		eviResult, err := CreateEvidence(ctx, db, cs, input)
		require.NoError(t, err)

		var evidenceID int64
		err = db.Get(&evidenceID, sq.Select("id").From("evidence").Where(sq.Eq{"uuid": eviResult.UUID}))
		require.NoError(t, err)
		op.Evidence = append(op.Evidence, models.Evidence{
			ID:          evidenceID,
			UUID:        eviResult.UUID,
			OperationID: id,
			OperatorID:  op.UserID,
			Description: eviResult.Description,
			ContentType: input.ContentType,
			OccurredAt:  eviResult.OccurredAt,
		})
	}
	makeEvidence(&goodOp, goodOp.ID, "item1")
	makeEvidence(&goodOp, goodOp.ID, "item2")
	makeEvidence(&goodOp, goodOp.ID, "item3")
	makeEvidence(&goodOp, goodOp.ID, "item4")

	makeEvidence(&badOp, badOp.ID, "item5")

	makeFinding := func(op *mockOperation, title string) {
		input := CreateFindingInput{OperationSlug: op.Op.Slug, Category: "Product", Title: title, Description: "desc"}
		findingResult, err := CreateFinding(ctx, db, input)
		require.NoError(t, err)

		var categoryID int64
		err = db.Get(&categoryID, sq.Select("id").From("finding_categories").Where(sq.Eq{"category": commonFindingCategories[0]}))
		require.NoError(t, err)
		var findingID int64
		err = db.Get(&findingID, sq.Select("id").From("findings").Where(sq.Eq{"uuid": findingResult.UUID}))
		require.NoError(t, err)
		op.Findings = append(op.Findings, models.Finding{
			ID:          findingID,
			UUID:        findingResult.UUID,
			Title:       findingResult.Title,
			Description: findingResult.Description,
			CategoryID:  &categoryID,
			OperationID: op.ID,
		})
	}

	makeFinding(&goodOp, "finding1")
	makeFinding(&goodOp, "finding2")
	makeFinding(&badOp, "finding3")

	return goodOp, badOp
}
