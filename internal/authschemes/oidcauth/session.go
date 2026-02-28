package oidcauth

import "encoding/gob"

// preLoginAuthSession is saved as authscheme session data before being redirected to okta
// so it can be verified on the callback route after returning
type preLoginAuthSession struct {
	Nonce              string
	StateChallengeCSRF string
	LoginMode          string
	OIDCService        string
}

// authSession is saved as authscheme session data after successfully authenticating as an okta user
type authSession struct {
	IdToken     string
	AccessToken string
}

func init() {
	gob.Register(&preLoginAuthSession{})
	gob.Register(&authSession{})
}
