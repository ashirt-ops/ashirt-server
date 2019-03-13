package dialog

import (
	"io"

	"github.com/manifoldco/promptui"
)

// UserQuery presents a free-answer dialog to the user with the given question (and default value,
// if one is provided). Returns the answer to the question, or an error if one is encountered.
func UserQuery(question string, defaultValue *string, inputStream io.ReadCloser) (string, error) {
	p := promptui.Prompt{
		Stdin:   inputStream,
		Label:   question,
		Pointer: promptui.PipeCursor,
	}
	if defaultValue != nil {
		p.Default = *defaultValue
	}
	return p.Run()
}

// MustUserQuery is a version of UserQuery that panics if UserQuery produces an error.
// This is typically unlikely, except during development.
func MustUserQuery(question string, defaultValue *string, inputStream io.ReadCloser) string {
	result, err := UserQuery(question, defaultValue, inputStream)
	if err != nil {
		panic(err)
	}
	return result
}
