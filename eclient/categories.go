package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// CategoryList is container struct for a list of CategoryResponses
type CategoryList struct {
	Object string              `json:"object"`
	Data   []*CategoryResponse `json:"data"`
}

// CategoryResponse contains the JSON response from the API.
type CategoryResponse struct {
	Object   string `json:"object"`
	ID       string `json:"id"`
	Segment  string `json:"segment"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	Lft      int    `json:"lft"`
	Rgt      int    `json:"rgt"`
	Depth    int    `json:"depth"`
	Created  time.Time
	Modified time.Time
}

// GetCategories returns a slice of categories.
func (c *EcomClient) GetCategories() ([]*CategoryResponse, error) {
	uri := c.endpoint + "/categories"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http new request failed")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.jwt)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "http do to %v failed", uri)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, errors.Wrapf(err, "client decode error")
		}
		return nil, errors.Errorf(fmt.Sprintf("Status: %d, Code: %s, Message: %s\n", e.Status, e.Code, e.Message))
	}

	var categoryList CategoryList
	if err := json.NewDecoder(res.Body).Decode(&categoryList); err != nil {
		return nil, errors.Wrapf(err, "json decode url %s failed", uri)
	}
	return categoryList.Data, nil
}
