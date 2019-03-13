// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package campaign

// Campaign stores stores the information about a specific campaign from the
// ASHIRT api server
type Campaign struct {
	Active      bool   `json:"active"`
	CreatedTime string `json:"createdTimestamp"`
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// New returns a pointer to a new Campaign
func New(id uint64, name string) *Campaign {
	return &Campaign{ID: id, Name: name}
}
