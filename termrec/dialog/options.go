package dialog

// Option is a small struct for providing options to a Select menu
type Option struct {
	Label       string
	Description string
	Action      func() OptionActionResponse
}

func (o Option) String() string {
	return o.Label
}

// Equals checks if two options are equal by comparing their labels
func (o Option) Equals(that Option) bool {
	return o.Label == that.Label
}

// OptionSliceToStringSlice conerts the given option slice into a string slice, for presenting
// to Select menus.
func OptionSliceToStringSlice(options []Option) []string {
	rtn := make([]string, len(options))
	for i, opt := range options {
		rtn[i] = opt.Label
	}
	return rtn
}

// OptionActionResponse contains a basic level of response information when selecting from a
// Select menu. This should help inform the parent of the select on how to proceed
type OptionActionResponse struct {
	// Err provides the error encountered when the option was selected, if any such error occurred.
	Err error

	// Value represents the single return value from that option action, if any
	Value interface{}

	// ShouldExit indicates if the option thinks that the user wants to leave this menu.
	ShouldExit bool
}

// NoAction is a shorthand OptionActionResponse that indicates that menu option handled everything
// A loose equivalent of a 201 response -- "it worked, but nothing to show"
func NoAction() OptionActionResponse {
	return OptionActionResponse{}
}

// PopAction is a shorthand OptionActionResponse that indicates the user wishes to leave this menu
// via returning to the previous menu. Analogous to `..` in unix filesystems
func PopAction() OptionActionResponse {
	return OptionActionResponse{ShouldExit: true}
}

// ChooseAction is a shorthand OptionActionResponse that indicates which _value_ / option the
// user selected.
func ChooseAction(val interface{}) func() OptionActionResponse {
	return func() OptionActionResponse {
		return OptionActionResponse{Value: val}
	}
}

// ErroredAction is a shorthand OptionActionResponse that indicates the option action encountered
// an error.
func ErroredAction(err error) OptionActionResponse {
	return OptionActionResponse{Err: err}
}

// MenuOptionGoBack provides a function that's directly interactable with Option. This communicates
// the desire to go up a level / exit the menu
func MenuOptionGoBack() OptionActionResponse {
	return PopAction()
}

// MenuContains checks if the given Option slice contains the passed option
func MenuContains(menu []Option, option Option) bool {
	for _, item := range menu {
		if option.Equals(item) {
			return true
		}
	}
	return false
}
