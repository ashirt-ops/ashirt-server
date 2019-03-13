// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package settings

import (
	"fmt"
	"os/user"
	"path"

	"github.com/theparanoids/ashirt/screenshotclient/config"
	"github.com/theparanoids/ashirt/screenshotclient/fuzzyfilefinder"
	"github.com/theparanoids/ashirt/screenshotclient/screencapture"
)

// Settings holds all the settings that modify the runtime behavior of the
// client
type Settings struct {
	Config         config.Config
	ProjectDefault string
	HomeDefault    string
}

// GetSettings returns the current settings based on the configuration file
func GetSettings() (*Settings, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	var s Settings

	c, err := config.Read()
	if err != nil {
		return nil, err
	}
	s.Config = c

	s.HomeDefault, err = screencapture.GetScreenshotDirectory()
	if err != nil {
		return nil, err
	}

	// if working folder doesn't exist use default of Ashirt on desktop
	if fuzzyfilefinder.Exists(c.WorkingFolder) {
		s.ProjectDefault = c.WorkingFolder
	} else {
		s.ProjectDefault = path.Join(u.HomeDir, "Desktop", "ashirt")

		fmt.Println("Warning: Working folder path : " + c.WorkingFolder + " Not Found using Default")
		fmt.Println("Make sure you are using the full path in the yaml config")
		fmt.Println("Default Working Directory" + s.ProjectDefault)
		fmt.Println()
	}

	return &s, nil
}
