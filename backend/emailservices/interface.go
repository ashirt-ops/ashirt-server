package emailservices

// EmailJob is a structure representing the required information needed to send a single email
// all other information should be specific to the email servicer, and should be provided to it
// in its configuration
type EmailJob struct {
	From        string
	To          string
	Subject     string
	Body        string
	OnCompleted func(error)
}

// EmailServicer is a simple interface for others to send emails. Once called, the expectation
// is that the email will eventually be sent, and called the OnCompleted when the email is sent,
// or fails to be sent
type EmailServicer interface {
	AddToQueue(EmailJob) error
}
