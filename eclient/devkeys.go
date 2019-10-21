package eclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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
		return nil, fmt.Errorf("%w: client: json marshal", err)
	}

	url := c.endpoint + "/developer-keys"
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("%w: client decode error", err)
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	var devKey DevKeyResponse
	if err = json.NewDecoder(res.Body).Decode(&devKey); err != nil {
		return nil, fmt.Errorf("%w: decode failed", err)
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
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var container DevKeysContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return container.Data, nil
}
