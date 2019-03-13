// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/theparanoids/ashirt/campaign"
	"github.com/theparanoids/ashirt/evidence"
	"github.com/theparanoids/ashirt/screenshotclient/config"
	"github.com/theparanoids/ashirt/shared"
)

var (
	// ErrBadRequest donotes a 400 response
	ErrBadRequest = errors.New("bad request")

	// ErrUnauthorized denotes a 401 response
	ErrUnauthorized = errors.New("unauthorized")

	// ErrNotFound denotes a 404 response
	ErrNotFound = errors.New("not found")

	// ErrInternalError denotes a 500 response
	ErrInternalError = errors.New("internal server error")

	// ErrUnknown is denotes an error that is not explicitly handled
	ErrUnknown = errors.New("unknown error")
)

// NewClient creates a new client
func NewClient(config *config.Config) *Client {
	return &Client{
		host:      config.APIURL,
		accessKey: config.AccessKey,
		secretKey: config.SecretKey,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Client is an HTTP client for communicating with the ASHIRT server
type Client struct {
	host      string
	client    http.Client
	accessKey string
	secretKey []byte
}

// GetCampaigns returns a list of campaigns
func (c *Client) GetCampaigns() ([]campaign.Campaign, error) {
	resp, err := c.Get("/api/operations")
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("unable to close response body: ", err)
		}
	}()

	err = checkStatus(resp.StatusCode, http.StatusOK)
	if err != nil {
		return nil, err
	}

	campaigns := make([]campaign.Campaign, 0)
	err = json.NewDecoder(resp.Body).Decode(&campaigns)

	return campaigns, err
}

// UploadScreenshot sends a screenshot to the ASHIRT API server to record as
// evidence
func (c *Client) UploadScreenshot(screenshot *evidence.Screenshot) error {
	body := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(body)

	fileWriter, err := bodyWriter.CreateFormFile("file", screenshot.FileName)
	if err != nil {
		return err
	}

	fp, err := os.Open(screenshot.FullPath)
	if err != nil {
		return err
	}
	defer func() {
		err := fp.Close()
		if err != nil {
			log.Println("unable to close file: ", err)
		}
	}()

	_, err = io.Copy(fileWriter, fp)
	if err != nil {
		return err
	}

	err = bodyWriter.WriteField("notes", screenshot.Description)
	if err != nil {
		return err
	}
	err = bodyWriter.WriteField("evidence_type", "1")
	if err != nil {
		return err
	}
	err = bodyWriter.WriteField("hash", screenshot.FileHash)
	if err != nil {
		return err
	}
	err = bodyWriter.WriteField("occurred_at", strconv.FormatInt(screenshot.OccurrenceTimestamp, 10))
	err = bodyWriter.Close()
	if err != nil {
		log.Println("unable to close body: ", err)
	}

	url := fmt.Sprintf("/api/operations/%d/evidence", screenshot.Campaign.ID)

	resp, err := c.Post(url, bodyWriter.FormDataContentType(), body)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("unable to close response body: ", err)
		}
	}()

	return checkStatus(resp.StatusCode, http.StatusCreated)
}

// Get creates a GET request and sends it to the server applying any additional processing needed
func (c *Client) Get(endpoint string) (*http.Response, error) {
	request, err := http.NewRequest("GET", c.host+endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(request)
}

// Post creates a POST request and sends it to the server applying any additional processing needed
func (c *Client) Post(endpoint, contentType string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("POST", c.host+endpoint, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)

	return c.Do(request)
}

// Do performs a request and provides any additional processing needed to communicate with the server
func (c *Client) Do(r *http.Request) (*http.Response, error) {
	r.Header.Set("Date", time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123))
	authorization, err := shared.BuildClientRequestAuthorization(r, c.accessKey, c.secretKey)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", authorization)

	return c.client.Do(r)
}

func checkStatus(status, expected int) error {
	if status == expected {
		return nil
	}

	switch status {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusInternalServerError:
		return ErrInternalError
	default:
		return ErrUnknown
	}
}
