// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package screencapture

import (
	"os/exec"
	"strings"
)

// GetScreenshotDirectory reads the screenshot directory from the macOS defaults
// system
func GetScreenshotDirectory() (string, error) {
	// runs this command defaults read com.apple.screencapture location
	cmdName := "defaults"
	cmdArgs := []string{"read", "com.apple.screencapture", "location"}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		return "", nil
	}

	sha := strings.Trim(string(cmdOut), "\n")
	return sha, nil
}

// SetScreenshotDirectory writes the screenshot directory to the macOS defaults
// system
func SetScreenshotDirectory(screenShotDir string) error {
	// runs this command defaults write com.apple.screencapture location /Users/someuser/Desktop/
	cmdName := "defaults"
	cmdArgs := []string{"write", "com.apple.screencapture", "location", screenShotDir}
	_, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		return err
	}

	return nil
}

// RemoveScreenshotDirectory will remove the location key from
// com.apple.screencapture for the case that the host didn't have one set.
// This brings the state back to the original rather than trying to put a
// false or bad value in or requiring the user to have the value set in order
// to use the program
func RemoveScreenshotDirectory() error {
	cmdName := "defaults"
	cmdArgs := []string{"delete", "com.apple.screencapture", "location"}
	_, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		return err
	}

	return nil
}

// ResetScreenshotDirectory intelligently changes the configured screenshot
// directory back to the specified one. This checks whether the directory is
// an empty string signifying that the no defaults was set and if so rather
// than setting the directory it removes the entry from defaults to bring
// the state back to the original
func ResetScreenshotDirectory(dir string) error {
	if dir == "" {
		return RemoveScreenshotDirectory()
	} else {
		return SetScreenshotDirectory(dir)
	}
}
