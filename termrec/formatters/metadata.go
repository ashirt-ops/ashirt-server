package formatters

import "fmt"

// Metadata allows for the capture of data for a particular recording session.
type Metadata struct {
	StartTimeUnix   int64
	DurationSeconds float64
	Title           string
	Shell           string
	Term            string
}

func (m Metadata) String() string {
	return fmt.Sprintf("%v %v %v %v %v", m.StartTimeUnix, m.DurationSeconds, m.Title, m.Shell, m.Term)
}
