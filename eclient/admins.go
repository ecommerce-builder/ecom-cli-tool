package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type createAdminRequest struct {
	Email  string `json:"email"`
	Passwd string `json:"password"`
	First  string `json:"firstname"`
	Last   string `json:"lastname"`
}

// CreateAdmin calls the API Service to create a new administrator.
func (c *EcomClient) CreateAdmin(email, passwd, first, last string) (*UserResponse, error) {
	p := createAdminRequest{
		Email:  email,
		Passwd: passwd,
		First:  first,
		Last:   last,
	}
	payload, err := json.Marshal(&p)
	if err != nil {
		return nil, fmt.Errorf("create admin failed: %w", err)
	}
	uri := c.endpoint + "/admins"
	res, err := c.request(http.MethodPost, uri, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP POST to %q return %s", uri, res.Status)
	}
	user := UserResponse{}
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("create product response decode failed: %w", err)
	}
	return &user, nil
}

// ListAdmins calls the API Service to get all administrators.
func (c *EcomClient) ListAdmins() ([]*UserResponse, error) {
	uri := c.endpoint + "/admins"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	users := make([]*UserResponse, 0, 8)
	err = json.NewDecoder(res.Body).Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("list users response decode failed: %w", err)
	}
	return users, nil
}

// DeleteAdmin calls the API service to delete an administrator with the
// given UUID.
func (c *EcomClient) DeleteAdmin(uuid string) error {
	uri := c.endpoint + "/admins/" + uuid
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()
	return nil
}
