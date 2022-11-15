// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

// BuildRequestHMAC builds a request HMAC from a secret key to authenticate a request for /api endpoints
// This function is shared by both client code to authenticate requests, and server to validate requests
//
// The return value is
//
//	base64(hmac-sha-256(VERB + "\n" +
//	                    REQUEST_PATH + "\n" +
//	                    DATE + "\n" +
//	                    sha256(REQUEST_BODY)
//	))
//
// It uses a separate requestBody argument instead of r.Body since reading from r.Body
// in both client & server will prevent reading the body again.
// Therefore it is the caller's responsibility to provide a separate request body reader.
// This is done by calling r.GetBody() on the client, and by reading request body to disk
// on server requests
func BuildRequestHMAC(r *http.Request, requestBody io.Reader, key []byte) []byte {
	requestBodySHA256 := sha256.New()
	io.Copy(requestBodySHA256, requestBody)

	m := new(bytes.Buffer)
	m.WriteString(r.Method)
	m.WriteString("\n")
	m.WriteString(r.URL.RequestURI())
	m.WriteString("\n")
	m.WriteString(r.Header.Get("Date"))
	m.WriteString("\n")
	m.Write(requestBodySHA256.Sum(nil))

	mac := hmac.New(sha256.New, key)
	mac.Write(m.Bytes())
	return mac.Sum(nil)
}

func BuildClientRequestAuthorization(r *http.Request, accessKey string, secretKey []byte) (string, error) {
	var body io.Reader
	if r.Method == "GET" {
		body = bytes.NewBuffer([]byte{})
	} else {
		var err error
		body, err = r.GetBody()
		if err != nil {
			return "", err
		}
	}

	hmac := BuildRequestHMAC(r, body, secretKey)
	return accessKey + ":" + base64.StdEncoding.EncodeToString(hmac), nil
}
