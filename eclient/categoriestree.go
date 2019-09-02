package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
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
		return errors.Wrapf(err, "client: json marshal failed")
	}

	uri := c.endpoint + "/categories-tree"
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPut, uri, body)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, "client decode error")
		}
		return errors.Errorf(fmt.Sprintf("Status: %d, Code: %s, Message: %s\n", e.Status, e.Code, e.Message))
	}
	return nil
}

// GetCategoriesTree returns the categories tree
func (c *EcomClient) GetCategoriesTree() (*CategoryTreeResponse, error) {
	uri := c.endpoint + "/categories-tree"
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

	var tree CategoryTreeResponse
	if err := json.NewDecoder(res.Body).Decode(&tree); err != nil {
		return nil, errors.Wrapf(err, "json decode url %s failed", uri)
	}
	return &tree, nil
}

// PurgeCatalog calls the API Service to purge the entire catalog.
func (c *EcomClient) PurgeCatalog() error {
	uri := c.endpoint + "/categories"
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		var e struct {
			Status  int    `json:"status"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, "client decode error")
		}
		return errors.Errorf(e.Message)
	}
	return nil
}
