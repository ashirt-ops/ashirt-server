package appdialogs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/theparanoids/ashirt/termrec/dialog"
	"github.com/theparanoids/ashirt/termrec/fancy"
	"github.com/theparanoids/ashirt/termrec/network"
	"github.com/pkg/errors"
)

var operationOptions = []dialog.Option{}

var menuOptionUpload = dialog.Option{Label: "Upload a file to the server", Action: showUploadSubmenu}
var menuOptionExit = dialog.Option{Label: "Exit", Action: dialog.MenuOptionGoBack}
var menuOptionUpdateOperations = dialog.Option{Label: "Refresh operations list", Action: updateOperationOptions}

// UserQuery is a re-packaging of dialog.UserQuery with inputStream pre-provided
func UserQuery(question string, defaultValue *string) (string, error) {
	return dialog.UserQuery(question, defaultValue, uploadStoreData.DialogInput)
}

// Select is a re-packaging of the dialog.Select with inpustStream pre-provided
func Select(label string, options []dialog.Option) dialog.OptionActionResponse {
	return dialog.Select(label, options, uploadStoreData.DialogInput)
}

// ShowUploadMainMenu presents the Main Menu during the uploading-of-capture phase
func ShowUploadMainMenu() {
	updateOperationOptions()

	mainMenuOptions := []dialog.Option{
		menuOptionUpdateOperations,
		menuOptionExit,
	}

	for {
		if !dialog.MenuContains(mainMenuOptions, menuOptionUpload) && len(operationOptions) > 0 {
			mainMenuOptions = append([]dialog.Option{menuOptionUpload}, mainMenuOptions...)
		}

		resp := Select("Select an operation", mainMenuOptions)
		if resp.ShouldExit {
			break
		}
		if resp.Err != nil {
			fmt.Println(fancy.Caution("Action failed", resp.Err))
		}
	}
}

func showUploadSubmenu() dialog.OptionActionResponse {
	err := tryUpload()
	if err != nil {
		switch errType := err.(type) {
		case CanceledOperation:
			SetDefaultData(errType.Data.(UploadDefaults))
			fmt.Println("Cancelled")
		default:
			fmt.Println("Encountered error during upload: " + err.Error())
		}
	}
	return dialog.NoAction()
}

func tryUpload() error {
	defaults := uploadStoreData.DefaultData
	if !network.BaseURLSet() {
		return errors.New("No service url specified -- check configuration")
	}

	path, err := UserQuery("Enter a filename", &defaults.FilePath)
	if err != nil {
		return errors.Wrap(err, "Could not retrieve filename")
	}

	fmt.Print("  Validating file... ")
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return errors.Wrap(err, "Unable to read recording")
	}
	fmt.Println(fancy.ClearLine(fancy.GreenCheck()+" File Validated", 0))

	slugResp := Select("Enter an Operation Slug", operationOptions)
	if slugResp.Err != nil {
		return errors.Wrap(slugResp.Err, "Could not retrieve operation slug")
	}

	description, err := UserQuery("Enter a description for this recording", &defaults.Description)
	if err != nil {
		return errors.Wrap(err, "Could not retrieve description")
	}

	// show a recap pre-upload
	fmt.Println(strings.Join([]string{
		fancy.WithBold("This will upload:", 0),
		fancy.WithBold("  File: ", 0) + fancy.WithBold(path, fancy.Yellow),
		fancy.WithBold("  Operation: ", 0) + fancy.WithBold(slugResp.Value.(string), fancy.Yellow),
		fancy.WithBold("  Description: ", 0) + fancy.WithBold(description, fancy.Yellow),
	}, "\n"))
	continueResp, err := dialog.YesNoPrompt("Do you want to continue?", "", uploadStoreData.DialogInput)
	if err != nil {
		return errors.Wrap(err, "Could not retrieve continue")
	}

	if continueResp == true {
		_, name := filepath.Split(path)

		input := network.UploadInput{
			OperationID: opIDFromSlug(slugResp.Value.(string)),
			Description: description,
			Filename:    name,
			Content:     bytes.NewReader(data),
		}

		var wg sync.WaitGroup
		wg.Add(1)
		stop := false
		var err error
		go func() {
			err = network.UploadToAshirt(input)
			wg.Done()
		}()
		go dialog.ShowLoadingAnimation("Loading", &stop)
		wg.Wait()
		stop = true
		fmt.Println(fancy.ClearLine(fancy.GreenCheck()+" File uploaded", 0))

		return errors.Wrap(err, "Could not upload")
	}

	return CanceledOperation{
		UploadDefaults{
			FilePath:      path,
			OperationSlug: slugResp.Value.(string),
			Description:   description,
		},
	}
}

func updateOperationOptions() (_ dialog.OptionActionResponse) {
	err := LoadOperations()
	if err != nil {
		fmt.Println(fancy.Caution("Unable to update operation list", err))
		return
	}

	operationOptions = make([]dialog.Option, len(uploadStoreData.Operations))
	for i, op := range uploadStoreData.Operations {
		operationOptions[i] = dialog.Option{
			Label:  op.Name,
			Action: dialog.ChooseAction(op.Slug),
		}
	}
	fmt.Printf("Loaded %v operations\n", len(uploadStoreData.Operations))
	return
}
