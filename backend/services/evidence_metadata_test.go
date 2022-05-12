// Copyright 2022, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"

	"github.com/stretchr/testify/require"
)

func TestCreateEvidenceMetadata(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	op := OpChamberOfSecrets
	evi := EviDobby

	input := services.EditEvidenceMetadataInput{
		OperationSlug: op.Slug,
		EvidenceUUID:  evi.UUID,
		Source:        "Insert Source",
		Body:          "some-body",
	}
	tryInsert := func(u models.User, input services.EditEvidenceMetadataInput) error {
		ctx := contextForUser(u, db)
		return services.CreateEvidenceMetadata(ctx, db, input)
	}

	// verify permissions
	require.Error(t, tryInsert(UserDraco, input))  // no operation access
	require.Error(t, tryInsert(UserSeamus, input)) // read access

	// verify insert
	require.NoError(t, tryInsert(UserRon, input)) // normal access

	metadataList := getEvidenceMetadataByEvidenceID(t, db, evi.ID)
	require.NotEmpty(t, metadataList)
	for _, v := range metadataList {
		if v.Source == input.Source {
			require.Equal(t, input.Body, v.Body)
			break
		}
	}
}

func TestUpdateEvidenceMetadata(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	op := OpChamberOfSecrets
	evi := EviDobby

	input := services.EditEvidenceMetadataInput{
		OperationSlug: op.Slug,
		EvidenceUUID:  evi.UUID,
		Source:        "Update Source",
		Body:          "some-body",
	}

	author := UserRon
	// update evidence metadata
	ctx := contextForUser(author, db)
	err := services.CreateEvidenceMetadata(ctx, db, input)
	require.NoError(t, err)

	tryUpdate := func(u models.User, input services.EditEvidenceMetadataInput) error {
		ctx := contextForUser(u, db)
		return services.UpdateEvidenceMetadata(ctx, db, input)
	}
	input.Body = "new-body"

	// verify permissions
	require.Error(t, tryUpdate(UserDraco, input))  // no operation access
	require.Error(t, tryUpdate(UserSeamus, input)) // read access

	// verify Update
	require.NoError(t, tryUpdate(author, input))

	metadataList := getEvidenceMetadataByEvidenceID(t, db, evi.ID)
	require.NotEmpty(t, metadataList)
	for _, v := range metadataList {
		if v.Source == input.Source {
			require.Equal(t, input.Body, v.Body)
			break
		}
	}
}

func TestUpsertEvidenceMetadata(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)
	op := OpChamberOfSecrets
	evi := EviDobby

	input := services.UpsertEvidenceMetadataInput{
		EditEvidenceMetadataInput: services.EditEvidenceMetadataInput{
			OperationSlug: op.Slug,
			EvidenceUUID:  evi.UUID,
			Source:        "Upsert Source",
			Body:          "",
		},
		Status:     "failed",
		Message:    helpers.Ptr("that didn't work"),
		CanProcess: helpers.Ptr(true),
	}

	tryUpsert := func(u models.User, input services.UpsertEvidenceMetadataInput) error {
		ctx := contextForUser(u, db)
		return services.UpsertEvidenceMetadata(ctx, db, input)
	}

	// verify permissions
	require.Error(t, tryUpsert(UserDraco, input))  // no operation access
	require.Error(t, tryUpsert(UserSeamus, input)) // read access

	author := UserRon
	// add evidence metadata
	require.NoError(t, tryUpsert(author, input))

	metadataList := getEvidenceMetadataByEvidenceID(t, db, evi.ID)
	require.NotEmpty(t, metadataList)
	_, metadataEntry := helpers.Find(metadataList, func(t models.EvidenceMetadata) bool {
		return t.Source == input.Source
	})
	require.NotNil(t, metadataEntry)
	require.Equal(t, input.Body, metadataEntry.Body)
	require.Equal(t, input.CanProcess, metadataEntry.CanProcess)
	require.Equal(t, input.Message, input.Message)
	require.Equal(t, input.Status, input.Status)

	input.Body = "new-body"
	input.Status = "sucess"
	input.Message = nil
	input.CanProcess = nil

	// verify Update
	require.NoError(t, tryUpsert(author, input))

	metadataList = getEvidenceMetadataByEvidenceID(t, db, evi.ID)
	require.NotEmpty(t, metadataList)
	_, metadataEntry = helpers.Find(metadataList, func(t models.EvidenceMetadata) bool {
		return t.Source == input.Source
	})
	require.NotNil(t, metadataEntry)
	require.Equal(t, input.Body, metadataEntry.Body)
	require.Equal(t, input.CanProcess, metadataEntry.CanProcess)
	require.Equal(t, input.Message, input.Message)
	require.Equal(t, input.Status, input.Status)
}
