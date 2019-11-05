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

// UserContainer object
type UserContainer struct {
	Object string          `json:"object"`
	Data   []*UserResponse `json:"data"`
	Links  interface{}     `json:"links,omitempty"`
}

// CreateUserRequest request body
type CreateUserRequest struct {
	Role      string `json:"role"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// UserResponse user JSON response body.
type UserResponse struct {
	Object      string    `json:"object"`
	ID          string    `json:"id"`
	UID         string    `json:"uid"`
	Role        string    `json:"role"`
	PriceListID string    `json:"price_list_id"`
	Email       string    `json:"email"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// CreateUser calls the API Service to create a new user.
func (c *EcomClient) CreateUser(ctx context.Context, u *CreateUserRequest) (*UserResponse, error) {
	request, err := json.Marshal(&u)
	if err != nil {
		return nil, fmt.Errorf("%w: client: json marshal", err)
	}

	uri := c.endpoint + "/users"
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, uri, body)
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

	var user UserResponse
	if err = json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("%w: decode failed", err)
	}
	return &user, nil
}

// GetUsers calls the API Service to retreieve a list of users.
func (c *EcomClient) GetUsers(ctx context.Context) ([]*UserResponse, error) {
	uri := c.endpoint + "/users"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: %w", res.Status, err)
	}

	var userContainer UserContainer
	if err := json.NewDecoder(res.Body).Decode(&userContainer); err != nil {
		return nil, errors.Wrapf(err, "json decode url", uri)
	}
	return userContainer.Data, nil
}
