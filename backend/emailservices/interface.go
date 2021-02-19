package emailservices

type EmailJob struct {
	From        string
	To          string
	Subject     string
	Body        string
	OnCompleted func(error)
}

type EmailServicer interface {
	AddToQueue(EmailJob) error
}
