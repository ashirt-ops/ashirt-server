package localauth

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/theparanoids/ashirt-server/backend"
)

// TOTPKey represents the secret and QR code for a given URL to authenticate with
type TOTPKey struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
	QRCode string `json:"qr"`
}

const totpDigits = otp.DigitsSix
const totpAlgorithm = otp.AlgorithmSHA1

func generateTOTP(accountName string) (*TOTPKey, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ashirt",
		AccountName: accountName,
		SecretSize:  64,
		Digits:      totpDigits,
		Algorithm:   totpAlgorithm,
	})
	if err != nil {
		return nil, err
	}
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)

	return &TOTPKey{
		URL:    key.URL(),
		Secret: key.Secret(),
		QRCode: "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}

func validateTOTP(passcode string, totpSecret string) error {
	if passcode == "" {
		return backend.MissingValueErr("TOTP Passcode")
	}
	if totpSecret == "" {
		return backend.MissingValueErr("TOTP Secret")
	}

	isValid, err := totp.ValidateCustom(passcode, totpSecret, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    totpDigits,
		Algorithm: totpAlgorithm,
	})
	if err != nil {
		return backend.InvalidTOTPErr(err)
	}
	if !isValid {
		return backend.InvalidTOTPErr(errors.New("totp.Validate failure"))
	}
	return nil
}
