package emailtemplates_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/database/seeding"
	"github.com/theparanoids/ashirt-server/backend/emailtemplates"
)

func TestBuildEmailContent(t *testing.T) {
	db := setupDb(t)

	allTemplates := []emailtemplates.EmailTemplate{
		emailtemplates.EmailRecoveryTemplate,
		emailtemplates.EmailRecoveryDeniedTemplate,
	}

	for _, tmpl := range allTemplates {
		body, subject, err := emailtemplates.BuildEmailContent(tmpl, emailtemplates.EmailTemplateData{
			DB:         db,
			UserRecord: &seeding.UserHarry,
		})
		require.NoError(t, err)
		require.NotEmpty(t, body)
		require.NotEmpty(t, subject)
	}
}

func setupDb(t *testing.T) *database.Connection {
	db := seeding.InitTestWithName(t, "emailtemplates-test-db")
	seeding.ApplySeeding(t, seeding.HarryPotterSeedData, db)

	return db
}
