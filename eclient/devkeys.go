package eclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ErrDeveloperKeyNotFound error
var ErrDeveloperKeyNotFound = errors.New("developer key not found")

// DevKeyRequest request body
type DevKeyRequest struct {
	UserID string `json:"user_id"`
}

// DevKeysContainer object
type DevKeysContainer struct {
	Object string            `json:"object"`
	Data   []*DevKeyResponse `json:"data"`
}

// DevKeyResponse developer key JSON response body.
type DevKeyResponse struct {
	Object   string    `json:"object"`
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Key      string    `json:"key"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// CreateDeveloperKey calls the API to create a developer key.
func (c *EcomClient) CreateDeveloperKey(ctx context.Context, d *DevKeyRequest) (*DevKeyResponse, error) {
	request, err := json.Marshal(&d)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	url := c.endpoint + "/developer-keys"
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w: request", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var devKey DevKeyResponse
	if err = json.NewDecoder(res.Body).Decode(&devKey); err != nil {
		return nil, fmt.Errorf("%w: decode", err)
	}
	return &devKey, nil
}

// GetDeveloperKeys returns a list of all developer keys for a given user.
func (c *EcomClient) GetDeveloperKeys(ctx context.Context, userID string) ([]*DevKeyResponse, error) {
	v := url.Values{}
	v.Set("user_id", userID)

	url := url.URL{
		Scheme:   "https",
		Host:     c.hostname,
		Path:     "developer-keys",
		RawQuery: v.Encode(),
	}
	res, err := c.request(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	var container DevKeysContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return container.Data, nil
}

// DeleteDeveloperKey calls the API service to attempt to delete a developer
// key by id.
func (c *EcomClient) DeleteDeveloperKey(ctx context.Context, devKeyID string) error {
	url := fmt.Sprintf("%s/developer-keys/%s", c.endpoint, devKeyID)
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrapf(err, "request failed url=%q", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return errors.Wrap(err, "decode")
		}
		if e.Code == "developer-keys/developer-key-not-found" {
			return ErrDeveloperKeyNotFound
		}
		return fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}
	return nil
}
