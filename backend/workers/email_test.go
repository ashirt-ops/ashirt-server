package workers_test

import (
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/database/seeding"
	"github.com/ashirt-ops/ashirt-server/backend/emailservices"
	"github.com/ashirt-ops/ashirt-server/backend/emailtemplates"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/workers"
	"github.com/stretchr/testify/require"
)

func setupDb(t *testing.T) *database.Connection {
	db := seeding.InitTestWithName(t, "workers-test-db")
	seeding.ApplySeeding(t, seeding.HarryPotterSeedData, db)

	return db
}

func TestEmailWorkerStartAndStop(t *testing.T) {
	db := setupDb(t)
	mailer := emailservices.MakeMemoryMailer(logging.NewNopLogger())

	emailWorker := workers.MakeEmailWorker(db, &mailer, logging.NewNopLogger())

	require.False(t, emailWorker.IsRunning())

	emailWorker.Start()

	require.True(t, emailWorker.IsRunning())

	emailWorker.Stop()

	require.False(t, emailWorker.IsRunning())
	require.Equal(t, 0, len(mailer.Outbox))
}

func TestEmailWorkerProcessEmail(t *testing.T) {
	db := setupDb(t)
	doneCh := make(chan bool)
	mailer := emailservices.MakeMemoryMailer(logging.NewNopLogger())
	emailWorker := workers.MakeEmailWorker(db, &mailer, logging.NewNopLogger())
	emailWorker.SleepAfterNoWorkDuration = 10 * time.Millisecond
	emailWorker.SleepAfterWorkDuration = 10 * time.Millisecond
	emailWorker.OnPassComplete = func() {
		doneCh <- true
	}

	targetUser := seeding.UserHarry
	givenTemplate := emailtemplates.EmailRecoveryTemplate // we need some valid template to produce data, so here we are
	badTemplate := "some-email-template"

	emailContent, err := emailtemplates.BuildEmailContent(givenTemplate, emailtemplates.EmailTemplateData{
		DB:         db,
		UserRecord: &targetUser,
	})
	require.NoError(t, err)

	db.Insert("email_queue", map[string]interface{}{
		"to_email": targetUser.Email,
		"user_id":  targetUser.ID,
		"template": badTemplate,
	})
	db.Insert("email_queue", map[string]interface{}{
		"to_email": targetUser.Email,
		"user_id":  targetUser.ID,
		"template": givenTemplate,
	})

	emailWorker.Start()

	<-doneCh
	emailWorker.Stop()
	sentEmails := mailer.Outbox[targetUser.Email]

	// check that the success email went through
	require.Equal(t, 1, len(sentEmails))
	resultEmail := sentEmails[0]
	require.Equal(t, emailContent.Subject, resultEmail.Subject)
	require.NotEmpty(t, resultEmail.Body) // recovery code won't match, so making the test a bit easier
	require.Equal(t, targetUser.Email, resultEmail.To)

	// check that the worker properly marked failures
	var failedEmail models.QueuedEmail
	err = db.Get(&failedEmail, sq.Select("*").From("email_queue").Where(sq.Eq{"template": badTemplate}))
	require.NoError(t, err)
	require.NotEmpty(t, failedEmail.ErrorText)
	require.GreaterOrEqual(t, int64(1), failedEmail.ErrorCount)
}
