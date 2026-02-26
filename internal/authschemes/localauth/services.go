package localauth

import (
	"context"
	"errors"
	"strings"

	"github.com/ashirt-ops/ashirt-server/internal/authschemes"
	"github.com/ashirt-ops/ashirt-server/internal/errorwrap"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/server/middleware"
	"golang.org/x/crypto/bcrypt"

	"github.com/ashirt-ops/ashirt-server/internal/services"
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
		return false, errorwrap.WrapError("Unable to check totp status for user", errorwrap.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanCheckTotp{UserID: userID}); err != nil {
		return false, errorwrap.WrapError("Non-Admin tried to check totp status for another user", errorwrap.UnauthorizedReadErr(err))
	}

	authData, err := bridge.FindUserAuthByUserID(userID)
	if err != nil {
		return false, errorwrap.WrapError("Unable to find auth details for user", errorwrap.UnauthorizedReadErr(err))
	}

	return (authData.TOTPSecret != nil), nil
}

func deleteUserTotp(ctx context.Context, bridge authschemes.AShirtAuthBridge, userSlug string) error {
	userID, err := services.SelfOrSlugToUserID(ctx, bridge.GetDatabase(), userSlug)
	if err != nil {
		return errorwrap.WrapError("Unable to delete totp for user", errorwrap.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanDeleteTotp{UserID: userID}); err != nil {
		return errorwrap.WrapError("Non-Admin tried to delete totp status for another user", errorwrap.UnauthorizedWriteErr(err))
	}
	authData, err := bridge.FindUserAuthByUserID(userID)

	if authData.TOTPSecret == nil {
		return errorwrap.BadInputErr(
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
		return errorwrap.WrapError("Unable to generate encrypted password", err)
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
