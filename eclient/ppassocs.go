package eclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ErrPPAssocNotFound error
var ErrPPAssocNotFound = errors.New("product to product association not found")

// PPAssocsContainer structure
type PPAssocsContainer struct {
	Object string     `json:"object"`
	Data   []*PPAssoc `json:"data"`
}

// PPAssoc represents a single product to product association.
type PPAssoc struct {
	Object         string    `json:"object"`
	ID             string    `json:"id"`
	PPAssocGroupID string    `json:"pp_assoc_group_id"`
	ProductFromID  string    `json:"product_from_id"`
	ProductToID    string    `json:"product_to_id"`
	Created        time.Time `json:"created"`
	Modified       time.Time `json:"modifed"`
}

// GetPPAssocs returns a list of all product to product associations
// for a given pp assoc group.
func (c *EcomClient) GetPPAssocs(ctx context.Context, ppaGroupID string) ([]*PPAssoc, error) {
	v := url.Values{}
	v.Set("pp_assoc_group_id", ppaGroupID)

	url := url.URL{
		Scheme:   "https",
		Host:     c.hostname,
		Path:     "products-assocs",
		RawQuery: v.Encode(),
	}
	res, err := c.request(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	var container PPAssocsContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return container.Data, nil
}

// DeletePPAssoc calls the API service to attempt to delete a product to product
// association
func (c *EcomClient) DeletePPAssoc(ctx context.Context, ppAssocID string) error {
	url := fmt.Sprintf("%s/products-assocs/%s", c.endpoint, ppAssocID)
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		return nil
	}
	if res.StatusCode == 400 {
		return ErrBadRequest
	}
	if res.StatusCode == 404 {
		return ErrPPAssocNotFound
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("err response %d", res.StatusCode)
	}
	return nil
}
