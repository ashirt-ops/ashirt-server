package localauth

import "encoding/gob"

// needsPasswordResetAuthSession is saved as an authscheme session for a user who has successfully
// authenticated with a username/password but requires a password reset.
type needsPasswordResetAuthSession struct {
	UserKey string
}

type needsTotpAuthSession struct {
	UserKey string
	Validated bool
}

func init() {
	gob.Register(&needsPasswordResetAuthSession{})
	gob.Register(&needsTotpAuthSession{})
}
