package dialog

import (
	"io"

	"github.com/manifoldco/promptui"
)

// MkBasicSelect provides a base for any Select operation. This essentially
// ensures that the given Select struct will read input from the proper source
func MkBasicSelect(inputStream io.ReadCloser) promptui.Select {
	return promptui.Select{
		Stdin:             inputStream,
		StartInSearchMode: false,
	}
}

// Select constructs a Select with the given properties:
//
// 1. Label and Options as given
// 2. SearcherContainsCI based on given labels
// 3. Interpreting of response to provide result
func Select(label string, options []Option, inputStream io.ReadCloser) OptionActionResponse {
	p := MkBasicSelect(inputStream)
	p.Label = label
	p.Searcher = SearcherContainsCI(OptionSliceToStringSlice(options))

	items := make([]string, len(options))
	for i, o := range options {
		items[i] = o.String()
	}
	p.Items = items

	selectedItem, _, err := p.Run()
	if err != nil {
		return ErroredAction(err)
	}
	return options[selectedItem].Action()
}
