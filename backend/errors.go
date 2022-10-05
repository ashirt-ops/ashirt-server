// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package backend

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	// ErrorDeprecated is an error that indicates a feature is deprecated. This would normally not
	// be used directly, but instead used to verify that a returned error is of "type" ErrorDeprecated
	ErrorDeprecated error = errors.New("warning: deprecated")
)

// HTTPError is a structure for communicating access/availability errors using a common format.
// Typically, users should opt for a pre-created error, rather than generate their own error.
//
// Note: all of these errors are designed to be communicated back to the API user.
type HTTPError struct {
	HTTPStatus   int
	PublicReason string
	WrappedError error
}

func (e *HTTPError) Error() string {
	return e.WrappedError.Error()
}

// Rewrap allows re-wrapping the wrapped error to include more information. The new error will
// consist of the newly provided message, followed by a colon, followed by the original error.
// e.g. err := HttpError(500, "private err", "public err"); err.Rewrap("outer"); // produces a wrapped error of "outer : private err"
func (e *HTTPError) Rewrap(msg string) {
	e.WrappedError = fmt.Errorf("%v : %w", msg, e.WrappedError)
}

// WrapError provides a mechanism to wrap any error in a consistent manner.
// If the error is an HTTPError (as defined in this package), then HTTPError.Rewrap will be called
// If the error is a regular error (i.e. non HTTPError), then a similar wrapping will occur, but the error will remain a non-HTTP error
// If the error is neither an HTTPError or regular error (i.e. the error is nil), then a new error will be generate with the provided message
func WrapError(msg string, err error) error {
	switch err := err.(type) {
	case *HTTPError:
		err.Rewrap(msg)
		return err
	case error:
		return fmt.Errorf("%v : %w", msg, err)
	}
	return fmt.Errorf(msg)
}

func HTTPErr(statusCode int, reason string, wrappedError error) error {
	return &HTTPError{
		HTTPStatus:   statusCode,
		PublicReason: reason,
		WrappedError: wrappedError,
	}
}

// BadInputErr provides a constructable, generic error for any request that, during use, does not make sense. Wraps a Bad Request error
func BadInputErr(err error, reason string) error {
	return HTTPErr(http.StatusBadRequest, reason, err)
}

// ServerErr provides a generic error for any error during a request, not covered by a more specific error
func ServerErr(err error) error {
	return HTTPErr(http.StatusInternalServerError, "Internal service error", err)
}

// SuggestiveServerErr provides an error with a customized message. This should be used primarily when you need to communicate
// how to fix an error
func SuggestiveServerErr(helpfulMessage string, err error) error {
	if helpfulMessage == "" {
		return ServerErr(err)
	}
	return HTTPErr(http.StatusInternalServerError, helpfulMessage, err)
}

// DatabaseErr provides a generic error for any database access error during a request
func DatabaseErr(err error) error {
	return HTTPErr(http.StatusInternalServerError, "Internal service error", err)
}

// SuggestiveDatabaseErr produces a 500 error, but sets the error message as something hopefully helpful to the user.
// this should be used sparingly, as it provides hints at what data is in the database.
func SuggestiveDatabaseErr(helpfulMessage string, err error) error {
	return HTTPErr(http.StatusInternalServerError, helpfulMessage, err)
}

// UploadErr provides an error for issues encountered while writing to the store
func UploadErr(err error) error {
	return HTTPErr(http.StatusInternalServerError, "The upload action could not be completed. Please try again.", err)
}

// DeleteErr provides an error for issues encountered while writing to the store
func DeleteErr(err error) error {
	return HTTPErr(http.StatusInternalServerError, "The delete action could not be completed. Please try again.", err)
}

// NotFoundErr provides an error for situations when a user requests data that does not exist.
func NotFoundErr(err error) error { return HTTPErr(http.StatusNotFound, "Not Found", err) }

// UnauthorizedReadErr provides an error for sitatutions where a user is unable to read whatever data is/may be found
func UnauthorizedReadErr(err error) error {
	return HTTPErr(http.StatusNotFound, "Not Found", err)
}

// UnauthorizedWriteErr provides an error for sitatutions where a user is unable to write/update data
func UnauthorizedWriteErr(err error) error {
	return HTTPErr(http.StatusUnauthorized, "Unauthorized", err)
}

// CSRFErr provides an error for when the CSRF validation fails
func CSRFErr(err error) error { return HTTPErr(http.StatusForbidden, "CSRF Failure", err) }

// BadAuthErr provides an error for sitatutions when a user authentication cannot be determined (mostly for alternative identity providers)
func BadAuthErr(err error) error { return HTTPErr(http.StatusForbidden, "Forbidden", err) }

// UserRequiresAdditionalAuthenticationErr is a helper for authschemes that need to redirect a user to a custom handler component
// on the frontend after a login attempt
func UserRequiresAdditionalAuthenticationErr(reason string) error {
	return HTTPErr(http.StatusPreconditionFailed, reason,
		//lint:ignore ST1005 Returned directly to the frontend for render
		fmt.Errorf("User requires additional auth: %s", reason),
	)
}

// InvalidPasswordErr provides an error for users that supply the wrong password.
// This wraps an Unauthorized status code
//
// Note: This should only be used when a user's existance is known -- e.g. when a user is trying to do some
// operation on their own account / admin doing some action on behalf of a user.
func InvalidPasswordErr(err error) error {
	return HTTPErr(http.StatusUnauthorized, "Invalid password", err)
}

// InvalidCredentialsErr provides an error for users that supply the wrong credentials.
// This wraps an Unauthorized status code
//
// Note: This should be used when a user's existance is _unknown_ -- i.e. the user is trying to login
func InvalidCredentialsErr(err error) error {
	return HTTPErr(http.StatusUnauthorized, "Invalid username or password", err)
}

// InvalidRecoveryErr provides an error for users that use an expired recovery code, or an incorrect
// recovery code.
// This wraps an Unauthorized status code
func InvalidRecoveryErr(err error) error {
	return HTTPErr(http.StatusUnauthorized, "Recovery code is invalid. Please ask an administrator to generate a new code.", err)
}

// MissingValueErr returns an error stating that some expected  value was not present.
// This is an alias for a Bad Request - type error
func MissingValueErr(valueName string) error {
	reason := fmt.Sprintf("Missing required field: %s", valueName)
	return HTTPErr(http.StatusBadRequest, reason, errors.New(reason))
}

// AccountDisabled returns an error indicating that a user's account has been disabled, and they
// cannot log in.
func AccountDisabled() error {
	err := DisabledUserError()
	return HTTPErr(http.StatusForbidden, err.Error(), err)
}

// IsErrorAccountDisabled checks if the provided error is the same as an "Account Disabled" error.
// See AccountDisabled() in this package.
func IsErrorAccountDisabled(err error) bool {
	switch err := err.(type) {
	case *HTTPError:
		model := AccountDisabled().(*HTTPError)
		return model.HTTPStatus == err.HTTPStatus && model.PublicReason == err.PublicReason
	}
	return false
}

// DisabledUserError is a version of AccountDisabled that returns an error, rather than an API Error
func DisabledUserError() error {
	//lint:ignore ST1005 Returned directly to the frontend for render
	return errors.New("This account has been disabled. Please contact an adminstrator if you think this is an error.")
}

// PanicedError represents any error the occurs
func PanicedError() error {
	return HTTPErr(http.StatusInternalServerError, "An unknown error occurred", errors.New("pancied during processing"))
}

// InvalidTOTPErr provides an error for users that provide an invalid TOTP passcode
func InvalidTOTPErr(err error) error {
	return HTTPErr(http.StatusUnauthorized, "Invalid passcode provided", err)
}

// FirstError returns the first non-nil error, or nil if all errors are nil.
// equivalement to:
//
//	if errs[0] != nil {return errs[0]} else if errs[1] != nil {return errs[1]} /*...*/ else { return nil }
func FirstError(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

// DeprecationWarning generates a wrapped error, with the given message and the underlying error as
// ErrorDeprecated
func DeprecationWarning(message string) error {
	return WrapError(message, ErrorDeprecated)
}

func WebauthnLoginError(err error, logMessage ...string) error {
	var fullErr error
	if len(logMessage) > 0 {
		fullErr = WrapError(strings.Join(logMessage, " ; "), err)
	} else {
		fullErr = err
	}

	return HTTPErr(http.StatusUnauthorized, "Unable to login. Verify your email", fullErr)
}
