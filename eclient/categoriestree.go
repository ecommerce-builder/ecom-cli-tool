package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// CategoryRequest body for updating the categories tree.
type CategoryRequest struct {
	Segment    string             `json:"segment"`
	Name       string             `json:"name"`
	Categories []*CategoryRequest `json:"categories"`
}

// CategoriesContainerResponse a container for list of catgory objects
type CategoriesContainerResponse struct {
	Object string                  `json:"object"`
	Data   []*CategoryTreeResponse `json:"data"`
}

// CategoryTreeResponse is a single node in the catalog.
type CategoryTreeResponse struct {
	Object     string                       `json:"object"`
	ID         string                       `json:"id"`
	Segment    string                       `json:"segment"`
	Name       string                       `json:"name"`
	Categories *CategoriesContainerResponse `json:"categories"`
	Products   *ProductContainerResponse    `json:"products,omitempty"`
}

// UpdateCategoriesTree calls the API Service to update the categories tree.
func (c *EcomClient) UpdateCategoriesTree(cats *CategoryRequest) error {
	request, err := json.Marshal(&cats)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	url := c.endpoint + "/categories-tree"
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPut, url, body)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		return fmt.Errorf("Status: %d, Code: %s, Message: %s: %w", e.Status, e.Code, e.Message, err)
	}
	return nil
}

// GetCategoriesTree returns the categories tree
func (c *EcomClient) GetCategoriesTree() (*CategoryTreeResponse, error) {
	uri := c.endpoint + "/categories-tree"
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

	var tree CategoryTreeResponse
	if err := json.NewDecoder(res.Body).Decode(&tree); err != nil {
		return nil, fmt.Errorf("json decode url %s failed: %w", uri, err)
	}
	return &tree, nil
}

// PurgeCatalog calls the API Service to purge the entire catalog.
func (c *EcomClient) PurgeCatalog() error {
	uri := c.endpoint + "/categories"
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		var e struct {
			Status  int    `json:"status"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("client decode error: %w", err)
		}
		return fmt.Errorf("%s: %w", e.Message, err)
	}
	return nil
}
