package emailservices

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddHtmlEmailToQueue(t *testing.T) {
	emailJob := EmailJob{
		To:       "Harry.Potter@hogwarts.edu",
		From:     "Hogwarts Admin<albus.dumbledore@hogwarts.edu>",
		Subject:  "HTML Body Test",
		HTMLBody: `<!DOCTYPE html><html><body>html body</body></html>`,
		Body:     "Plaintext body",
		OnCompleted: func(err error) {
			require.NoError(t, err)
		},
	}

	content, err := buildEmailContent(emailJob)

	require.NoError(t, err)
	strContent := string(content)

	fmt.Println("Message: =======\n", strContent, "\n===========")

	expectedHeaders := []string{
		fmt.Sprintf("To: %v\r\n", emailJob.To),
		fmt.Sprintf("From: %v\r\n", emailJob.From),
		fmt.Sprintf("Subject: %v\r\n", emailJob.Subject),
		"MIME-Version: 1.0\r\n",
	}
	preamble := strings.Join(expectedHeaders, "")
	actualEmailPreamble := strContent[:len(preamble)]
	require.Equal(t, preamble, actualEmailPreamble)
	mergedParts := strContent[len(preamble):]

	partLines := strings.Split(mergedParts, "\r\n")
	mpContentTypeLine := partLines[0]
	boundaryText := "boundary="
	boundaryIndex := strings.Index(mpContentTypeLine, boundaryText)
	boundary := mpContentTypeLine[boundaryIndex+len(boundaryText):]

	expectedContentLines := []string{
		fmt.Sprintf("Content-Type: multipart/alternative; boundary=%v\r\n", boundary),
		"\r\n",
		fmt.Sprintf("--%v\r\n", boundary),
		"Content-Type: text/plain; charset=\"utf-8\"\r\n",
		"\r\n",
		emailJob.Body,
		"\r\n\r\n", // an extra line break after the main body
		fmt.Sprintf("--%v\r\n", boundary),
		"Content-Type: text/html; charset=\"utf-8\"\r\n",
		"\r\n",
		emailJob.HTMLBody,
		"\r\n\r\n", // an extra line break after the main body
		fmt.Sprintf("--%v--\r\n", boundary),
	}
	expectedContent := strings.Join(expectedContentLines, "")
	require.Equal(t, expectedContent, mergedParts)
}
