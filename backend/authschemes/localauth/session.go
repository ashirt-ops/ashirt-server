// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package localauth

import (
	"encoding/gob"
	"net/http"

	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
)

// localAuthSession is saved as an authscheme session for users that have "some difficulty" in logging in --
// i.e. a plain authentication is insufficient, and more action is required. Speciifically, this
// comes in the following flavors:
//   - User must reset their password
//   - User must supply their TOTP code
type localAuthSession struct {
	SessionValid  bool
	Username      string
	TOTPValidated bool
}

func init() {
	gob.Register(&localAuthSession{})
}

func readLocalSession(r *http.Request, bridge authschemes.AShirtAuthBridge) *localAuthSession {
	sess, ok := bridge.ReadAuthSchemeSession(r).(*localAuthSession)
	if !ok {
		return &localAuthSession{SessionValid: false}
	}
	return sess
}

func (sess *localAuthSession) writeLocalSession(w http.ResponseWriter, r *http.Request, bridge authschemes.AShirtAuthBridge) {
	bridge.SetAuthSchemeSession(w, r, &localAuthSession{
		SessionValid:  true,
		Username:      sess.Username,
		TOTPValidated: sess.TOTPValidated,
	})
}
