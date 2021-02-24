package emailservices_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/emailservices"
	"github.com/theparanoids/ashirt-server/backend/logging"
)

func TestAddToQueue(t *testing.T) {
	writer := bytes.NewBuffer(make([]byte, 0, 1024))
	servicer := emailservices.MakeWriterMailer(writer, logging.NewNopLogger())
	onCompletedCalled := false
	servicer.AddToQueue(emailservices.EmailJob{
		To:      "Harry.Potter@hogwarts.edu",
		From:    "Hogwarts Admin<albus.dumbledore@hogwarts.edu>",
		Subject: "Tuition",
		Body:    "Gimme my money!",
		OnCompleted: func(err error) {
			onCompletedCalled = true
			require.NoError(t, err)
		},
	})
	content := string(writer.Bytes())

	expectedOutput :=
		`
======================================================
To:      | Harry.Potter@hogwarts.edu
From:    | Hogwarts Admin<albus.dumbledore@hogwarts.edu>
Subject: | Tuition
Body:
Gimme my money!
======================================================
`

	require.Equal(t, content, expectedOutput)
	require.True(t, onCompletedCalled)
}
