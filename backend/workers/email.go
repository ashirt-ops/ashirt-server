package workers

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/emailservices"
	"github.com/theparanoids/ashirt-server/backend/emailtemplates"
	"github.com/theparanoids/ashirt-server/backend/logging"
)

// EmailStatus reflects the possible statuses for emails in any state
type EmailStatus = string

const (
	// Created reflects emails where the email has been acknowledged, but not yet sent
	Created EmailStatus = "created"
	// Sent reflects emails that have been attempted to be delivered (and probably succeeded)
	Sent EmailStatus = "sent"
	// Errored reflects emails that could not be sent for any reason
	Errored EmailStatus = "error"
)

type EmailWorker struct {
	db         *database.Connection
	emailQueue chan emailservices.EmailJob
	stopChan   chan bool
	running    bool
	servicer   emailservices.EmailServicer
	logger     logging.Logger
}

func MakeEmailWorker(db *database.Connection, servicer emailservices.EmailServicer, logger logging.Logger) EmailWorker {
	emailCh := make(chan emailservices.EmailJob)
	stopCh := make(chan bool)
	return EmailWorker{
		db:         db,
		emailQueue: emailCh,
		stopChan:   stopCh,
		servicer:   servicer,
		logger:     logger,
	}
}

func (w *EmailWorker) GetEmailQueue() *chan emailservices.EmailJob {
	return &w.emailQueue
}

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

func (w *EmailWorker) Stop() {
	w.stopChan <- true
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
			Where(sq.Eq{"email_status": []string{Created, Errored}}).
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
			sleepDuration = 60
		} else {
			sleepDuration = 20
		}

		time.Sleep(sleepDuration * time.Second)
	}
}

func (w *EmailWorker) queueEmail(email emailRequest) error {
	user, err := w.db.RetrieveUserByID(email.UserID)

	if err != nil {
		return err
	}
	templateData := emailtemplates.EmailTemplateData{
		SendToAddress: email.To,
		UserRecord:    &user,
		DB:            w.db,
	}
	body, subject, err := emailtemplates.BuildEmailContent(email.Template, templateData)

	if err != nil {
		return err
	}
	if w.servicer == nil {
		return fmt.Errorf("Email servicer has not been assigned to email worker")
	}
	err = w.servicer.AddToQueue(emailservices.EmailJob{
		Body:    body,
		Subject: subject,
		To:      email.To,
		From:    config.EmailFromAddress(),
		OnCompleted: func(encounteredErr error) {
			if encounteredErr != nil {
				w.logger.Log("msg", "Unable to send email", "err", encounteredErr.Error())
				w.db.Update(sq.Update("email_queue").
					Set("error_count", sq.Expr("error_count + 1")).
					SetMap(map[string]interface{}{
						"email_status": Errored,
						"error_text":   encounteredErr.Error(),
					}).
					Where(sq.Eq{"id": email.EmailID}))
			} else {
				err := w.db.Update(sq.Update("email_queue").
					Set("email_status", Sent).
					Where(sq.Eq{"id": email.EmailID}))
				if err != nil {
					w.logger.Log("msg", "Unable to set email completed status", "error", err.Error())
				}
			}
		},
	})
	return err
}
