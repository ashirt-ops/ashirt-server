// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package config

import (
	"strings"
)

// Supported Flags
const (
	// WelcomeFlag is for testing purposes -- displays a welcome message on the operations page
	WelcomeFlag string = "welcome-message"
	// AllowMetadataEdit tells the frontend that it can render the metadata-editing capabilities
	AllowMetadataEdit string = "allow-metadata-edit"
)

var flags []string

// Flags returns a list of all of the flags that were loaded from the environment.
// This is cached for speedier access later
func Flags() []string {
	if len(flags) == 0 {
		flags = strings.Split(app.Flags, ",")
	}
	return flags
}

func HasFlag(flagName string) bool {
	allFlags := Flags()
	for _, f := range allFlags { // TODO: replace with FindMatch (in other branch)
		if f == flagName {
			return true
		}
	}
	return false
}
