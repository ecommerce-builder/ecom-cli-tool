package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ListUsers calls the API Service to retreieve a list of users.
func (c *EcomClient) ListUsers() ([]*User, error) {
	uri := c.endpoint + "/customers"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: %w", res.Status, err)
	}

	var users []*User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("json decode url %s failed: %w", uri, err)
	}
	return users, nil
}
