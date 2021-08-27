package emailservices

// EmailJob is a structure representing the required information needed to send a single email
// all other information should be specific to the email servicer, and should be provided to it
// in its configuration
type EmailJob struct {
	From        string
	To          string
	Subject     string
	Body        string
	HTMLBody    string
	OnCompleted func(error)
}

// EmailServicer is a simple interface for others to send emails. Once called, the expectation
// is that the email will eventually be sent, and called the OnCompleted when the email is sent,
// or fails to be sent
type EmailServicer interface {
	AddToQueue(EmailJob) error
}

// EmailServicerType acts as an enum for selecting known email servicers types
type EmailServicerType = string

const (
	// StdOutEmailer refers to an email servicer that simply outputs the emails to the terminal.
	// Useful for local testing
	StdOutEmailer EmailServicerType = "stdout"

	// MemoryEmailer refers to an email servicer that holds all emails sent in memory.
	// Useful for unit testing
	MemoryEmailer EmailServicerType = "memory"

	// SMTPEmailer refers to an email servicer that sends emails via SMTP.
	SMTPEmailer EmailServicerType = "smtp"
)
