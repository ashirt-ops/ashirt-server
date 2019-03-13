// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/theparanoids/ashirt/campaign"
	"github.com/theparanoids/ashirt/screenshotclient/settings"

	"github.com/martinlindhe/inputbox"
)

// MaxFileName is the maximum length allowed for a filename. This is
// realistically going to be different from OS to OS and FS to FS but 255
// seems to be a pretty safe number for now and we can lower this if needed.
// We should probably make this a fixed number, even if a specific OS os FS
// can go larger to ensure that behavior is consistent across OS and FS
const MaxFileName = 255

var (
	// ErrCancelled is returned when a user cancels specifying a description
	ErrCancelled         = errors.New("screenshot cancelled")
	previousDescriptions string
	// This can be expanded but is probably fine for now
	filenameReplace = regexp.MustCompile(`[^0-9a-zA-Z\-]+`)
)

// Screenshot holds the metadata related to a screenshot to be sent to the
// ASHIRT API server for recording
type Screenshot struct {
	ID                  int
	FileName            string
	FullPath            string
	Description         string
	FileHash            string
	OccurrenceTimestamp int64
	Campaign            campaign.Campaign
}

// GetFileHash returns the sha256 hash of the file
func GetFileHash(fileName string) (string, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", os.ErrNotExist
	}
	hasher := sha256.New()
	s, err := ioutil.ReadFile(fileName)
	hasher.Write(s)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// ProcessFile cleans up and organizes evidence
func ProcessFile(screenshot *Screenshot, settings *settings.Settings) error {
	//remove spaces
	screenshot.Campaign.Name = strings.Replace(screenshot.Campaign.Name, " ", "", -1)
	currentTime := time.Now().Local()

	if _, err := os.Stat(settings.ProjectDefault); os.IsNotExist(err) {
		err = os.Mkdir(settings.ProjectDefault, os.ModePerm)
		if err != nil {
			return err
		}
	}
	// the project path with the campaign
	projectPathCampaign := path.Join(settings.ProjectDefault, screenshot.Campaign.Name)

	if _, err := os.Stat(projectPathCampaign); os.IsNotExist(err) {
		err = os.Mkdir(projectPathCampaign, os.ModePerm)
		if err != nil {
			return err
		}
	}
	// the project path for the day the screenshot was tkaen
	projectPathDay := path.Join(settings.ProjectDefault, screenshot.Campaign.Name, currentTime.Format("2006-01-02"))

	if _, err := os.Stat(projectPathDay); os.IsNotExist(err) {
		err = os.Mkdir(projectPathDay, os.ModePerm)
		if err != nil {
			return err
		}
	}
	//set the file name to the time
	//replace : with .
	timeNow := strings.Replace(time.Now().Format(time.RFC3339), ":", ".", -1)
	// Compute the length to truncate of the description for the filename
	descLen := MaxFileName - (len(timeNow) + 1 + len(".png"))
	//replace whitespace with underscore
	description := substr(filenameReplace.ReplaceAllString(screenshot.Description, "_"), descLen)
	filename := timeNow + "_" + description + ".png"

	fullRename := filepath.Join(projectPathDay, filename)

	// Move file to the archive folder
	err := os.Rename(screenshot.FullPath, fullRename)
	if err != nil {
		return err
	}

	screenshot.FullPath = fullRename

	return nil
}

// Exists checks to see if a file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CollectScreenshot checks that a screenshot has a png extension and requests a
// description for it
func CollectScreenshot(projectDir string, fileName string) (Screenshot, error) {
	var screenshot Screenshot
	screenshot.FileName = fileName
	screenshot.FullPath = path.Join(projectDir, fileName)

	//prompt user
	response, ok := inputbox.InputBox("A-Shirt Capture", "Enter a description", previousDescriptions)
	if !ok {
		return screenshot, ErrCancelled
	}

	previousDescriptions = response
	screenshot.Description = response
	screenshot.OccurrenceTimestamp = time.Now().UnixNano()
	return screenshot, nil
}

// substr is a rune aware substring function that returns a string that is at
// most len bytes wide to the closest complete rune in the original string
func substr(s string, length int) string {
	var sub = make([]byte, 0)

	for _, r := range s {
		bytes := len(sub)
		if bytes+utf8.RuneLen(r) > length {
			break
		}

		sub = append(sub, []byte(string(r))...)
	}

	return string(sub)
}
