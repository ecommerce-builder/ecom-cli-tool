package eclient

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// A Category is a single node in the catalog.
type Category struct {
	Segment    string     `json:"segment" yaml:"segment"`
	Name       string     `json:"name" yaml:"name"`
	Categories []Category `json:"categories" yaml:"categories"`
}

// A Catalog contains a single root node of the catalog.
type Catalog struct {
	Endpoints []string `yaml:"endpoints"`
	Category  Category `yaml:"catalog"`
}

// UpdateCatalog calls the API Service to update the catalog.
func (c *EcomClient) UpdateCatalog(cats Category) error {
	payload, err := json.Marshal(&cats)
	if err != nil {
		return errors.Wrapf(err, "client: json marshal failed")
	}
	uri := c.endpoint + "/catalog"
	body := strings.NewReader(string(payload))
	res, err := c.request(http.MethodPut, uri, body)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return errors.Errorf("HTTP PUT to %q return %s", uri, res.Status)
	}
	return nil
}

// PurgeCatalog calls the API Service to purge the entire catalog.
func (c *EcomClient) PurgeCatalog() error {
	uri := c.endpoint + "/catalog"
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		var e struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, "client decode error")
		}
		return errors.Errorf(e.Message)
	}
	return nil
}
