package fancy

import "strings"

// see https://misc.flogisoft.com/bash/tip_colors_and_formatting

const (
	Bold int = 1 << (iota)
	Dim
	Underlined
	Blink
	Reverse
	Hidden
	Black
	Red
	Green
	BrownOrange
	Blue
	Purple
	Cyan
	LightGray
	DarkGray
	LightRed
	LightGreen
	Yellow
	LightBlue
	LightPurple
	LightCyan
	White
)

func flagMap() map[int]string {
	return map[int]string{
		Bold:        sBold,
		Dim:         sDim,
		Underlined:  sUnderlined,
		Blink:       sBlink,
		Reverse:     sReverse,
		Hidden:      sHidden,
		Black:       cBlack,
		Red:         cRed,
		Green:       cGreen,
		BrownOrange: cBrownOrange,
		Blue:        cBlue,
		Purple:      cPurple,
		Cyan:        cCyan,
		LightGray:   cLightGray,
		DarkGray:    cDarkGray,
		LightRed:    cLightRed,
		LightGreen:  cLightGreen,
		Yellow:      cYellow,
		LightBlue:   cLightBlue,
		LightPurple: cLightPurple,
		LightCyan:   cLightCyan,
		White:       cWhite,
	}
}

const escape = "\033["

// Plain resets the coloring to the default terminal value
const Plain = "\033[0m"

// Clear removes all text from the line
const Clear = "\033[2K"

const (
	sBold       string = "1"
	sDim        string = "2"
	sUnderlined string = "4"
	sBlink      string = "5"
	sReverse    string = "7"
	sHidden     string = "8"
)

const (
	cBlack       string = "30"
	cRed         string = "31"
	cGreen       string = "32"
	cBrownOrange string = "33"
	cBlue        string = "34"
	cPurple      string = "35"
	cCyan        string = "36"
	cLightGray   string = "37"
	cDarkGray    string = "90"
	cLightRed    string = "91"
	cLightGreen  string = "92"
	cYellow      string = "93"
	cLightBlue   string = "94"
	cLightPurple string = "95"
	cLightCyan   string = "96"
	cWhite       string = "97"
	cClear       string = "2K"
)

// WithPizzazz allows for text to be rendered with terminal styling (colors, weight, underline + other effects)
func WithPizzazz(s string, flags int) string {
	prefix := ""
	codes := []string{}
	for k, v := range flagMap() {
		if flags&k > 0 {
			codes = append(codes, v)
		}
	}
	if len(codes) > 0 {
		prefix = escape + strings.Join(codes, ";") + "m"
	}

	return prefix + s + Plain
}

// WithBold is a shorthand for WithPizzazz(s, flags | Bold)
func WithBold(s string, otherFlags ...int) string {
	combinedFlags := Bold | mergeFlags(otherFlags)
	return WithPizzazz(s, combinedFlags)
}

// ClearLine removes all text from the current line, resets the cursor to the start of the line,
// then calls WithPizzazz(s, flags)
func ClearLine(s string, flags ...int) string {
	return WithPizzazz(Clear+"\r"+s, mergeFlags(flags))
}

// GreenCheck renders a green unicode checkmark
func GreenCheck() string {
	return WithPizzazz("✔", LightGreen)
}

// RedCross renders a red unicode cross / X mark
func RedCross() string {
	return WithPizzazz("✘", Red)
}

// Caution renders a message in yellow (and red) indicating that some issue occurred
func Caution(message string, err error) string {
	cautionMsg := WithPizzazz("! ", Bold|Yellow) + WithPizzazz(message, Yellow)
	if err != nil {
		cautionMsg += " : " + WithPizzazz(err.Error(), Red)
	}
	return cautionMsg
}

func Fatal(message string, err error) string {
	errMsg := RedCross() + " " + WithPizzazz(message, Bold|Red)

	if err != nil {
		errMsg += " : " + WithPizzazz(err.Error(), Red)
	}
	return errMsg
}

func mergeFlags(flags []int) int {
	rtn := 0
	for _, val := range flags {
		rtn |= val
	}
	return rtn
}
