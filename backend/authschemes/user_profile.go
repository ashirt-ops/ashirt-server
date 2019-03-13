// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package authschemes

import (
	"github.com/theparanoids/ashirt/backend/services"
)

// UserProfile containes the necessary information to create a new user
type UserProfile struct {
	FirstName string
	LastName  string
	Slug      string
	Email     string
}

// ToCreateUserInput converts the given UserProfile into a more useful services.CreateUserInput
func (up UserProfile) ToCreateUserInput() services.CreateUserInput {
	return services.CreateUserInput{
		FirstName: up.FirstName,
		LastName:  up.LastName,
		Slug:      up.Slug,
		Email:     up.Email,
	}
}
