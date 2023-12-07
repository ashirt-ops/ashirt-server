package workers

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ashirt-ops/ashirt-server/backend/config"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/emailservices"
	"github.com/ashirt-ops/ashirt-server/backend/emailtemplates"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
)

// EmailStatus reflects the possible statuses for emails in any state
type EmailStatus = string

const (
	// EmailCreated reflects emails where the email has been acknowledged, but not yet sent
	EmailCreated EmailStatus = "created"
	// EmailSent reflects emails that have been attempted to be delivered (and probably succeeded)
	EmailSent EmailStatus = "sent"
	// EmailErrored reflects emails that could not be sent for any reason
	EmailErrored EmailStatus = "error"
)

// EmailWorker is a struct that creates the functionality of reading emails from the email queue (db table)
// and passing those to the email servicer
type EmailWorker struct {
	db                       *database.Connection
	emailQueue               chan emailservices.EmailJob
	stopChan                 chan bool
	running                  bool
	servicer                 emailservices.EmailServicer
	logger                   logging.Logger
	SleepAfterWorkDuration   time.Duration
	SleepAfterNoWorkDuration time.Duration
	OnPassComplete           func()
}

// MakeEmailWorker constructs an EmailWorker
func MakeEmailWorker(db *database.Connection, servicer emailservices.EmailServicer, logger logging.Logger) EmailWorker {
	emailCh := make(chan emailservices.EmailJob)
	stopCh := make(chan bool)
	return EmailWorker{
		db:                       db,
		emailQueue:               emailCh,
		stopChan:                 stopCh,
		servicer:                 servicer,
		logger:                   logger,
		SleepAfterWorkDuration:   20 * 1000,
		SleepAfterNoWorkDuration: 60 * 1000,
	}
}

// GetEmailQueue returns the channel that the EmailWorker uses to add emails to the outgoing queue.
// The email servicer should read from this channel to know what emails to send out.
func (w *EmailWorker) GetEmailQueue() *chan emailservices.EmailJob {
	return &w.emailQueue
}

// Start starts the email worker's processing. Note that calling this while the worker is already
// running will do nothing
func (w *EmailWorker) Start() {
	if !w.running {
		w.running = true
		defer func() {
			if r := recover(); r != nil {
				w.logger.Log("msg", "recovered from worker panic", "error", r)
			}
		}()
		w.start()
	}
}

// start _actually_ starts the worker.
func (w *EmailWorker) start() {
	w.logger.Log("msg", "Starting worker")
	go w.run()
	go func() {
		<-w.stopChan
		w.running = false
	}()
}

// Stop stops the email worker at it's next opporunity (generally between individual email sends)
// any email that has been picked up by the servicer will run to completion.
func (w *EmailWorker) Stop() {
	w.stopChan <- true
}

// IsRunning returns the current state of the worker. If running, the worker is processing new emails
// returns true if the worker is running, false otherwise.
func (w *EmailWorker) IsRunning() bool {
	return w.running
}

type emailRequest struct {
	EmailID  int64  `db:"id"`
	To       string `db:"to_email"`
	UserID   int64  `db:"user_id"`
	Template string `db:"template"`
}

func (w *EmailWorker) run() {
	var sleepDuration time.Duration = 60
	for w.running {
		var emails []emailRequest

		// get emails from email queue
		err := w.db.Select(&emails, sq.Select("id", "to_email", "user_id", "template").
			From("email_queue").
			Where(sq.Eq{"email_status": []string{EmailCreated, EmailErrored}}).
			Where(sq.Expr("error_count < ?", 3)).
			OrderBy("updated_at ASC"). // grab the oldest emails first, prefer jobs that have not errored out
			Limit(50))

		if len(emails) > 0 {
			for _, email := range emails {
				if !w.running {
					break
				}
				err = w.queueEmail(email)
				if err != nil {
					w.logger.Log("msg", "Unable to queue email", "error", err.Error())
					continue
				}
			}
			sleepDuration = w.SleepAfterWorkDuration
		} else {
			sleepDuration = w.SleepAfterNoWorkDuration
		}

		if w.OnPassComplete != nil {
			w.OnPassComplete()
		}
		time.Sleep(sleepDuration * time.Millisecond)
	}
}

func (w *EmailWorker) queueEmail(email emailRequest) error {
	user, err := w.db.RetrieveUserByID(email.UserID)

	if err != nil {
		return err
	}
	templateData := emailtemplates.EmailTemplateData{
		UserRecord: &user,
		DB:         w.db,
	}
	emailContent, err := emailtemplates.BuildEmailContent(email.Template, templateData)

	if err != nil {
		setEmailFailed(w.db, email.EmailID, w.logger, err)
		return err
	}
	if w.servicer == nil {
		return fmt.Errorf("Email servicer has not been assigned to email worker")
	}
	err = w.servicer.AddToQueue(emailservices.EmailJob{
		Body:     emailContent.PlaintTextContent,
		HTMLBody: emailContent.HTMLContent,
		Subject:  emailContent.Subject,
		To:       email.To,
		From:     config.EmailFromAddress(),
		OnCompleted: func(encounteredErr error) {
			if encounteredErr != nil {
				setEmailFailed(w.db, email.EmailID, w.logger, encounteredErr)
			} else {
				err := w.db.Update(sq.Update("email_queue").
					Set("email_status", EmailSent).
					Where(sq.Eq{"id": email.EmailID}))
				if err != nil {
					w.logger.Log("msg", "Unable to set email completed status", "error", err.Error())
				}
			}
		},
	})
	return err
}

func setEmailFailed(db *database.Connection, emailID int64, logger logging.Logger, err error) {
	logger.Log("msg", "Unable to send email", "err", err.Error())
	db.Update(sq.Update("email_queue").
		Set("error_count", sq.Expr("error_count + 1")).
		SetMap(map[string]interface{}{
			"email_status": EmailErrored,
			"error_text":   err.Error(),
		}).
		Where(sq.Eq{"id": emailID}))
}
