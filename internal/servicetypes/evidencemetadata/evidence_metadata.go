package evidencemetadata

import "github.com/ashirt-ops/ashirt-server/internal/helpers"

type Status string

const (
	// StatusQueued reflects the state where a worker has received the task, but decided it could not process immediately.
	StatusQueued Status = "Queued"
	// StatusProcessing reflects the state where a task has been given to the worker, but the worker hasn't responded with the result yet.
	StatusProcessing Status = "Processing"
	// StatusCompleted reflects a work-done scenario where the processing succeeded
	// (note: this does not mean that we have data -- just that no error occurred.)
	StatusCompleted Status = "Completed"
	// StatusCompleted reflects a work-done scenario where the processing failed, and an error was returned
	StatusError Status = "Error"
)

func (v Status) String() string {
	return string(v)
}

func (v Status) Ptr() *Status {
	return helpers.Ptr(v)
}
