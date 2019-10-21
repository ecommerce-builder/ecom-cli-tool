package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CategoryContainer is container struct for a list of Category
type CategoryContainer struct {
	Object string      `json:"object"`
	Data   []*Category `json:"data"`
}

// Category contains the JSON response from the API.
type Category struct {
	Object   string    `json:"object"`
	ID       string    `json:"id"`
	Segment  string    `json:"segment"`
	Path     string    `json:"path"`
	Name     string    `json:"name"`
	Lft      int       `json:"lft"`
	Rgt      int       `json:"rgt"`
	Depth    int       `json:"depth"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// GetCategories returns a slice of categories.
func (c *EcomClient) GetCategories() ([]*Category, error) {
	uri := c.endpoint + "/categories"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("http new request failed: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.jwt)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do to %v failed: %w", uri, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("client decode error: %w", err)
		}
		return nil, fmt.Errorf("Status: %d, Code: %s, Message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	var categoryContainer CategoryContainer
	if err := json.NewDecoder(res.Body).Decode(&categoryContainer); err != nil {
		return nil, fmt.Errorf("json decode url %s failed: %w", uri, err)
	}
	return categoryContainer.Data, nil
}
