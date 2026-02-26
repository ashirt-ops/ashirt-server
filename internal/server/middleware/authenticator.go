package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/errors"
	"github.com/ashirt-ops/ashirt-server/internal/logging"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/server/remux"
	"github.com/ashirt-ops/ashirt-server/internal/session"

	sq "github.com/Masterminds/squirrel"
)

var policyCtxKey = &struct{ name string }{"policy"}
var userCtxKey = &struct{ name string }{"userID"}
var adminCtxKey = &struct{ name string }{"admin"}

// InjectPolicy is a helper function to add a policy to the context under the expected key
func InjectPolicy(ctx context.Context, p policy.Policy) context.Context {
	return context.WithValue(ctx, policyCtxKey, p)
}

// InjectUser is a helper function to add a user to the context under the expected key
func InjectUser(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userCtxKey, userID)
}

// InjectAdmin is a helper function to add a flag denoting if the current user is an Administrator
func InjectAdmin(ctx context.Context, isAdmin bool) context.Context {
	return context.WithValue(ctx, adminCtxKey, isAdmin)
}

// InjectIntoContextInput is a small structure for defining what is placed into the context
type InjectIntoContextInput struct {
	IsSuperAdmin bool
	UserID       int64
	UserPolicy   policy.Policy
}

// InjectIntoContext adds a collection of data to the appropriate keys into the context.
func InjectIntoContext(ctx context.Context, i InjectIntoContextInput) context.Context {
	ctxValues := map[interface{}]interface{}{
		userCtxKey:   i.UserID,
		adminCtxKey:  i.IsSuperAdmin,
		policyCtxKey: i.UserPolicy,
	}
	for k, v := range ctxValues {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

// IsAdmin is used to check if the current user has been identified as an admin. Note that this
// value will only change when the session store is cleared for this user (i.e. they log out)
func IsAdmin(ctx context.Context) bool {
	admin, _ := ctx.Value(adminCtxKey).(bool)
	return admin
}

// UserID is used to retrieve the user id from context
func UserID(ctx context.Context) int64 {
	id, _ := ctx.Value(userCtxKey).(int64)
	return id
}

// Policy is used to retrieve policy from context
func Policy(ctx context.Context) policy.Policy {
	p, ok := ctx.Value(policyCtxKey).(policy.Policy)
	if !ok {
		return &policy.Deny{}
	}
	return p
}

func AuthenticateAppAndInjectCtx(db *database.Connection) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, cleanup, err := cloneBody(r)
			if err != nil {
				respondWithError(w, r, errors.WrapError("Unable to clone http body", err))
				return
			}
			defer cleanup()

			userData, err := authenticateAPI(db, r, body)
			if err != nil {
				logging.LogWithoutAuth(
					"Unable to build user policy",
					"error", err.Error(),
				)
				respondWithError(w, r, errors.UnauthorizedWriteErr(err))
				return
			}
			ctx := buildContextForUser(r.Context(), db, userData.ID, false, userData.Headless)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthenticateUserAndInjectCtx(db *database.Connection, sessionStore *session.Store) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := sessionStore.Read(r)
			if sess.UserID == 0 {
				next.ServeHTTP(w, r)
				return
			}
			// users that log in to the web (where this is used) cannot be headless users
			ctx := buildContextForUser(r.Context(), db, sess.UserID, sess.IsAdmin, false)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func BuildContextForUser(ctx context.Context, db *database.Connection, userID int64, isSuperAdmin, isHeadless bool) context.Context {
	return buildContextForUser(ctx, db, userID, isSuperAdmin, isHeadless)
}

func buildContextForUser(ctx context.Context, db *database.Connection, userID int64, isSuperAdmin, isHeadless bool) context.Context {
	return InjectIntoContext(ctx, InjectIntoContextInput{
		IsSuperAdmin: isSuperAdmin,
		UserID:       userID,
		UserPolicy:   buildPolicyForUser(ctx, db, userID, isSuperAdmin, isHeadless),
	})
}

func buildPolicyForUser(ctx context.Context, db *database.Connection, userID int64, isSuperAdmin, isHeadless bool) policy.Policy {
	var roles []models.UserOperationPermission

	var groupRoles []models.UserGroupOperationPermission

	err := db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.Select(&roles, sq.Select("operation_id", "role").
			From("user_operation_permissions").
			Where(sq.Eq{"user_id": userID}))

		var userGroupIds []int64
		tx.Select(&userGroupIds, sq.Select("group_id").
			From("group_user_map").
			Where(sq.Eq{"user_id": userID}))

		tx.Select(&groupRoles, sq.Select("operation_id", "role").
			From("user_group_operation_permissions").
			Where(sq.Eq{"group_id": userGroupIds}))
	})

	if err != nil {
		logging.ReqLogger(ctx).Error("Unable to build user policy", "error", err.Error())
		return &policy.Deny{}
	}
	roleMap := make(map[int64]policy.OperationRole)
	for _, role := range roles {
		roleMap[role.OperationID] = role.Role
	}
	for _, role := range groupRoles {
		val, ok := roleMap[role.OperationID]
		noRole := !ok
		assignedRoleIsLowest := ok && val == policy.OperationRoleRead
		groupRoleIsHigher := ok && val == policy.OperationRoleWrite && role.Role == policy.OperationRoleAdmin

		if noRole || assignedRoleIsLowest || groupRoleIsHigher {
			roleMap[role.OperationID] = role.Role
		}
	}
	return &policy.Union{
		P1: policy.NewAuthenticatedPolicy(userID, isSuperAdmin),
		P2: &policy.Operation{
			UserID:           userID,
			IsHeadless:       isHeadless,
			OperationRoleMap: roleMap,
		},
	}
}

func respondWithError(w http.ResponseWriter, r *http.Request, err error) {
	remux.JSONHandler(func(r *http.Request) (interface{}, error) { return nil, err }).ServeHTTP(w, r)
}

// cloneBody saves the request body of non-GET requests to a file on disk since the
// request body may be quite large (screenshot or other evidence) and we need to read
// it twice, once to validate the HMAC of the API request, and then the actual
// business logic needs to read it again. This is similar to what go does internally
// for ParseForm, but since HMAC needs to be validated for all request types including
// application/json we cannot rely on that.
//
// The returned `cleanup` function removes the temporary file created and should be
// called after the request is completed
func cloneBody(r *http.Request) (io.Reader, func(), error) {
	if r.Method == "GET" {
		return bytes.NewBuffer([]byte{}), func() {}, nil
	}
	bodyTmpFile, err := os.CreateTemp("", "ashirt-body")
	if err != nil {
		return nil, func() {}, err
	}

	_, err = io.Copy(bodyTmpFile, r.Body)
	if err != nil {
		return nil, func() {}, err
	}
	bodyTmpFile.Close()
	r.Body.Close()

	r.Body, err = os.Open(bodyTmpFile.Name())
	if err != nil {
		return nil, func() {}, err
	}

	body, err := os.Open(bodyTmpFile.Name())
	if err != nil {
		return nil, func() {}, err
	}

	cleanup := func() {
		body.Close()
		r.Body.Close()
		err := os.Remove(bodyTmpFile.Name())
		if err != nil {
			logging.LogWithoutAuth(
				"Unable to remove tmp file",
				"error", err,
				"tmpFile", bodyTmpFile.Name(),
			)
		}
	}
	return body, cleanup, nil
}
