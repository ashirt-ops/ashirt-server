package evidencemetadata

import "github.com/theparanoids/ashirt-server/backend/helpers"

type Status string

const (
	StatusUnaccepted Status = "Unaccepted"
	StatusError      Status = "Error"
	StatusQueued     Status = "Queued"
	StatusCompleted  Status = "Completed"
)

func (v Status) String() string {
	return string(v)
}

func (v Status) Ptr() *Status {
	return helpers.Ptr(v)
}
