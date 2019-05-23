package eclient

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// ListCustomers call the API Service to retreieve a list of customers.
func (c *EcomClient) ListCustomers() ([]*Customer, error) {
	uri := c.endpoint + "/customers"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, errors.Wrapf(err, "%s", res.Status)
	}

	var customers []*Customer
	if err := json.NewDecoder(res.Body).Decode(&customers); err != nil {
		return nil, errors.Wrapf(err, "json decode url %s failed", uri)
	}
	return customers, nil
}
