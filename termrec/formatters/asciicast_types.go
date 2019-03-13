package formatters

// ASCIICastHeader is a structure for representing Asciinema formatted files. See here:
// https://github.com/asciinema/asciinema/blob/develop/doc/asciicast-v2.md
// Note that this compatible with V2 explicitly. Future revisions are not guaranteed.
type ASCIICastHeader struct {
	// only first 3 are required
	Version       int               `json:"version"`
	Width         uint16            `json:"width"`
	Height        uint16            `json:"height"`
	Timestamp     int64             `json:"timestamp"`
	Duration      float64           `json:"duration,omitempty"`
	IdleTimeLimit float64           `json:"idle_time_limit,omitempty"`
	Command       string            `json:"command,omitempty"`
	Title         string            `json:"title,omitempty"`
	Env           map[string]string `json:"env"`
	Theme         *ASCIICastTheme   `json:"theme,omitempty"`
}

// ASCIICastTheme is the struct that represents the Asciinema file theme sub-structure
type ASCIICastTheme struct {
	Foreground string `json:"fg,omitempty"`
	Background string `json:"bg,omitempty"`
	Palette    string `json:"palette,omitempty"`
}

// ASCIInemaEvent is the struct that represents the Asciinema file event sub-structure
// Note: in the asciicast format, the data is represented as an array, so there is no
// explicit json field naming here. See (f ASCIICast) WriteEvent for how these get turned
// into the proper format
type ASCIInemaEvent struct {
	When float64
	Type string
	Data string
}
