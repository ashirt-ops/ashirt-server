package localauth

import (
	"context"
	"errors"
	"strings"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/authschemes"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"golang.org/x/crypto/bcrypt"

	"github.com/theparanoids/ashirt-server/backend/services"
)

type RegistrationInfo struct {
	Password           string
	Username           string
	Email              string
	FirstName          string
	LastName           string
	ForceResetPassword bool
}

func readUserTotpStatus(ctx context.Context, bridge authschemes.AShirtAuthBridge, userSlug string) (bool, error) {
	userID, err := services.SelfOrSlugToUserID(ctx, bridge.GetDatabase(), userSlug)
	if err != nil {
		return false, backend.WrapError("Unable to check totp status for user", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanCheckTotp{UserID: userID}); err != nil {
		return false, backend.WrapError("Non-Admin tried to check totp status for another user", backend.UnauthorizedReadErr(err))
	}

	authData, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return false, backend.WrapError("Unable to find auth details for user", backend.UnauthorizedReadErr(err))
	}

	return (authData.TOTPSecret != nil), nil
}

func deleteUserTotp(ctx context.Context, bridge authschemes.AShirtAuthBridge, userSlug string) error {
	userID, err := services.SelfOrSlugToUserID(ctx, bridge.GetDatabase(), userSlug)
	if err != nil {
		return backend.WrapError("Unable to delete totp for user", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanDeleteTotp{UserID: userID}); err != nil {
		return backend.WrapError("Non-Admin tried to delete totp status for another user", backend.UnauthorizedWriteErr(err))
	}
	authData, err := bridge.FindUserAuthByUserID(userID)

	if authData.TOTPSecret == nil {
		return backend.BadInputErr(
			errors.New("User does not have a TOTP key associated"),
			"This account does not have a TOTP key",
		)
	}

	authData.TOTPSecret = nil
	return bridge.UpdateAuthForUser(authData)
}

func registerNewUser(ctx context.Context, bridge authschemes.AShirtAuthBridge, info RegistrationInfo) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return backend.WrapError("Unable to generate encrypted password", err)
	}

	userResult, err := bridge.CreateNewUser(authschemes.UserProfile{
		FirstName: info.FirstName,
		LastName:  info.LastName,
		Slug:      strings.ToLower(info.FirstName + "." + info.LastName),
		Email:     info.Email,
	})
	if err != nil {
		return err
	}
	return bridge.CreateNewAuthForUser(authschemes.UserAuthData{
		UserID:             userResult.UserID,
		Username:           info.Username,
		EncryptedPassword:  encryptedPassword,
		NeedsPasswordReset: info.ForceResetPassword,
	})
}
