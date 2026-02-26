package localauth

import (
	"bytes"
	"encoding/base64"
	stderrors "errors"
	"image/png"
	"time"

	"github.com/ashirt-ops/ashirt-server/internal/errors"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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
		return errors.MissingValueErr("TOTP Passcode")
	}
	if totpSecret == "" {
		return errors.MissingValueErr("TOTP Secret")
	}

	isValid, err := totp.ValidateCustom(passcode, totpSecret, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    totpDigits,
		Algorithm: totpAlgorithm,
	})
	if err != nil {
		return errors.InvalidTOTPErr(err)
	}
	if !isValid {
		return errors.InvalidTOTPErr(stderrors.New("totp.Validate failure"))
	}
	return nil
}
