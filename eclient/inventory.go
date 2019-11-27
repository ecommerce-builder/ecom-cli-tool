package eclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// ErrInventoryNotFound error
var ErrInventoryNotFound = errors.New("inventory not found")

// UpdateInventoryRequest JSON payload
type UpdateInventoryRequest struct {
	Onhand      *int  `json:"onhand,omitempty"`
	Overselling *bool `json:"overselling,omitempty"`
}

// Inventory holds inventory for a single product
type Inventory struct {
	Object      string    `json:"object"`
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	ProductPath string    `json:"product_path"`
	ProductSKU  string    `json:"product_sku"`
	Onhand      int       `json:"onhand"`
	Overselling bool      `json:"overselling"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// InventoryContainer list of inventory objects.
type InventoryContainer struct {
	Object string       `json:"object"`
	Data   []*Inventory `json:"data"`
}

// InventoryBatchUpdateContainer list of inventory update request.
type InventoryBatchUpdateContainer struct {
	Object string                         `json:"object"`
	Data   []*InventoryBatchUpdateRequest `json:"data"`
}

//InventoryBatchUpdateRequest JSON batch inventory update request.
type InventoryBatchUpdateRequest struct {
	ProductID   *string `json:"product_id"`
	Onhand      *int    `json:"onhand"`
	Overselling *bool   `json:"overselling"`
}

// GetInventory calls the API service to retrieve an inventory
// by id.
func (c *EcomClient) GetInventory(ctx context.Context, invID string) (*Inventory, error) {
	url := fmt.Sprintf("%s/inventory/%s", c.endpoint, invID)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err,
			"request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		if e.Code == "inventory/inventory-not-found" {
			return nil, ErrInventoryNotFound
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}

	var v Inventory
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &v, nil
}

// GetAllInventory calls the API service to retrieve all inventory
// for all products.
func (c *EcomClient) GetAllInventory(ctx context.Context) ([]*Inventory, error) {
	url := fmt.Sprintf("%s/inventory", c.endpoint)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err,
			"request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	var container InventoryContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}
	return container.Data, nil
}

// UpdateInventory calls the API service to update an individual inventory.
func (c *EcomClient) UpdateInventory(ctx context.Context, invID string, req *UpdateInventoryRequest) (*Inventory, error) {
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("%w: client: json marshal failed", err)
	}

	body := bytes.NewReader(request)
	url := fmt.Sprintf("%s/inventory/%s", c.endpoint, invID)
	res, err := c.request(http.MethodPatch, url, body)
	if err != nil {
		return nil, errors.Wrapf(err,
			"request(http.MethodPatch, url=%q, nil)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "client decode")
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var v Inventory
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &v, nil
}

// UpdateInventoryBatch calls the API service to
func (c *EcomClient) UpdateInventoryBatch(ctx context.Context, inv []*InventoryBatchUpdateRequest) ([]*Inventory, error) {
	reqContainer := InventoryBatchUpdateContainer{
		Object: "list",
		Data:   inv,
	}

	b, err := json.Marshal(&reqContainer)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal")
	}
	body := bytes.NewReader(b)
	url := fmt.Sprintf("%s/inventory:batch-update", c.endpoint)
	res, err := c.request(http.MethodPatch, url, body)
	if err != nil {
		return nil, errors.Wrapf(err,
			"request(http.MethodPatch, url=%q, body=%q)", url, string(b))
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrapf(err, "client decode")
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var container InventoryContainer
	if err = json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return container.Data, nil
}
