package emailservices

import (
	"log/slog"
)

// MemoryMailer is an EmailServicer that holds all of the emails it receives in memory. This mailer
// is designed to be used with testing, where the caller can quickly check if the email details
// were received correctly
type MemoryMailer struct {
	Outbox map[string][]EmailJob // emails stored as To address : [emails, ...]
	logger *slog.Logger
}

// MakeMemoryMailer constructs a MemoryMailer
func MakeMemoryMailer(logger *slog.Logger) MemoryMailer {
	return MemoryMailer{
		Outbox: make(map[string][]EmailJob),
	}
}

// AddToQueue adds the given email job to memory (specifically, to the MemoryMailer.Outbox)
// this will never return an error, nor call OnCompleted with an error (but OnCompleted _will_ still
// be called)
func (m *MemoryMailer) AddToQueue(job EmailJob) error {
	if _, ok := m.Outbox[job.To]; !ok {
		m.Outbox[job.To] = make([]EmailJob, 0, 1)
	}
	m.Outbox[job.To] = append(m.Outbox[job.To], job)

	job.OnCompleted(nil)

	return nil
}
