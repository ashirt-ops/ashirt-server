package localauth

import "encoding/gob"

// needsPasswordResetAuthSession is saved as an authscheme session for a user who has successfully
// authenticated with a username/password but requires a password reset.
type needsPasswordResetAuthSession struct {
	UserKey string
}

func init() {
	gob.Register(&needsPasswordResetAuthSession{})
}
