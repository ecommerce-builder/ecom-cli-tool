package eclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ErrWebhookExists webhook already exists
var ErrWebhookExists = errors.New("webhook exists")

// ErrWebhookNotFound webhook not found error
var ErrWebhookNotFound = errors.New("webhook not found")

// ErrEventTypeNotFound error
var ErrEventTypeNotFound = errors.New("event type not found")

// WebhookEventsRequest events
type WebhookEventsRequest struct {
	Object string   `json:"object"`
	Data   []string `json:"data"`
}

// CreateWebhookRequest payload
type CreateWebhookRequest struct {
	URL    string               `json:"url"`
	Events WebhookEventsRequest `json:"events"`
}

// UpdateWebhookRequest payload
type UpdateWebhookRequest struct {
	URL     string               `json:"url"`
	Events  WebhookEventsRequest `json:"events"`
	Enabled bool                 `json:"enabled"`
}

// WebhookContainer structure
type WebhookContainer struct {
	Object string             `json:"object"`
	Data   []*WebhookResponse `json:"data"`
}

// WebhookResponse object
type WebhookResponse struct {
	Object     string    `json:"object"`
	ID         string    `json:"id"`
	SigningKey string    `json:"signing_key"`
	URL        string    `json:"url"`
	Events     []string  `json:"events"`
	Enabled    bool      `json:"enabled"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
}

// CreateWebhook calls the API to create a new webhook.
func (c *EcomClient) CreateWebhook(ctx context.Context, p *CreateWebhookRequest) (*WebhookResponse, error) {
	request, err := json.Marshal(&p)
	if err != nil {
		return nil, fmt.Errorf("%w: client: json marshal", err)
	}

	url := c.endpoint + "/webhooks"
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
			return nil, fmt.Errorf("%w: client decode", err)
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	var webhook WebhookResponse
	if err = json.NewDecoder(res.Body).Decode(&webhook); err != nil {
		return nil, fmt.Errorf("%w: decode", err)
	}
	return &webhook, nil
}

// GetWebhook calls the API service to retrive a single webhook.
func (c *EcomClient) GetWebhook(ctx context.Context, webhookID string) (*WebhookResponse, error) {
	url := c.endpoint + "/webhooks/" + webhookID
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		return nil, ErrBadRequest
	}
	if res.StatusCode == 404 {
		return nil, ErrWebhookNotFound
	}

	var w WebhookResponse
	if err := json.NewDecoder(res.Body).Decode(&w); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &w, nil
}

// GetWebhooks returns a list of all webhooks.
func (c *EcomClient) GetWebhooks(ctx context.Context) ([]*WebhookResponse, error) {
	uri := c.endpoint + "/webhooks"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var container WebhookContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return container.Data, nil
}

// UpdateWebhook calls the API service to update a webhook
func (c *EcomClient) UpdateWebhook(ctx context.Context, webhookID string, req *UpdateWebhookRequest) (*WebhookResponse, error) {
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	body := strings.NewReader(string(request))
	url := c.endpoint + "/webhooks/" + webhookID
	res, err := c.request(http.MethodPatch, url, body)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		return nil, fmt.Errorf("Status: %d, Code: %s, Message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	var w WebhookResponse
	if err := json.NewDecoder(res.Body).Decode(&w); err != nil {
		return nil, errors.Wrapf(err, "json decode url %s", url)
	}
	return &w, nil
}

// DeleteWebhook calls the API service to delete a webhook
func (c *EcomClient) DeleteWebhook(ctx context.Context, webhookID string) error {
	url := c.endpoint + "/webhooks/" + webhookID
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return ErrWebhookNotFound
	}

	return nil
}
