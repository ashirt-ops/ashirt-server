package dialog

import (
	"io"

	"github.com/theparanoids/ashirt/termrec/fancy"

	"github.com/manifoldco/promptui"
)

type booleanOption struct {
	Label         string
	Value         bool
	Icon          string
	ActiveStyle   func(interface{}) string
	SelectedStyle func(interface{}) string
}

var yes = booleanOption{
	Label:         " Yes ",
	Icon:          fancy.GreenCheck(),
	Value:         true,
	ActiveStyle:   promptui.Styler(promptui.BGGreen, promptui.FGBold, promptui.FGBlack),
	SelectedStyle: promptui.Styler(promptui.FGGreen),
}
var no = booleanOption{
	Label:         " No ",
	Icon:          fancy.RedCross(),
	Value:         false,
	ActiveStyle:   promptui.Styler(promptui.BGRed, promptui.FGBold, promptui.FGWhite),
	SelectedStyle: promptui.Styler(promptui.FGRed),
}

var yesNo = []booleanOption{yes, no}

// YesNoPrompt spans a "Select" dialog, where a given user will be propted with a
// yes/no question (plus, optional details, if details is not the empty string)
// will return (true, nil) if the user selected "Yes", (false, nil) if the user selected "No"
// (false, <error>) if some error occurred.
func YesNoPrompt(label, details string, inputStream io.ReadCloser) (bool, error) {
	yesNoAsStringSlice := make([]string, len(yesNo))
	for i := range yesNo {
		yesNoAsStringSlice[i] = yesNo[i].Label
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "{{ .Label | call .ActiveStyle }}",
		Inactive: "{{ .Label }}",
		Selected: "{{ (print .Icon (.Label | call .SelectedStyle)) }}",
	}
	if details != "" {
		templates.Details = details
	}
	p := MkBasicSelect(inputStream)
	p.Items = yesNo
	p.Label = label
	p.Searcher = SearcherContainsCI(yesNoAsStringSlice)
	p.Templates = templates

	selectedIndex, _, err := p.Run()
	return yesNo[selectedIndex].Value, err
}
