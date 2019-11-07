package eclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"net/http"

	"github.com/pkg/errors"
)

// ErrPriceListNotFound error
var ErrPriceListNotFound = errors.New("price list not found")

// ErrPriceListCodeExists error
var ErrPriceListCodeExists = errors.New("price list code already exists")

// PriceListContainer price list container JSON body.
type PriceListContainer struct {
	Object string       `json:"object"`
	Data   []*PriceList `json:"data"`
}

// PriceList price list JSON response body.
type PriceList struct {
	Object        string    `json:"object"`
	ID            string    `json:"id"`
	PriceListCode string    `json:"price_list_code"`
	CurrencyCode  string    `json:"currency_code"`
	Strategy      string    `json:"strategy"`
	IncTax        bool      `json:"inc_tax"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modified"`
}

// CreatePriceListRequest request body for creating a new price list.
type CreatePriceListRequest struct {
	PriceListCode string `json:"price_list_code"`
	CurrencyCode  string `json:"currency_code"`
	Strategy      string `json:"strategy"`
	IncTax        bool   `json:"inc_tax"`
	Name          string `json:"name"`
	Description   string `json:"description"`
}

// UpdatePriceListRequest request body for updating a new price list.
type UpdatePriceListRequest struct {
	PriceListCode string `json:"price_list_code,omitempty"`
	CurrencyCode  string `json:"currency_code,omitempty"`
	Strategy      string `json:"strategy,omitempty"`
	IncTax        bool   `json:"inc_tax,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
}

// CreatePriceList calls the API service to attempt to create a new price list.
func (c *EcomClient) CreatePriceList(ctx context.Context, req *CreatePriceListRequest) (*PriceList, error) {
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrapf(err, "json marshal")
	}
	url := fmt.Sprintf("%s/price-lists", c.endpoint)
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "request")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		if e.Code == "price-lists/price-list-not-found" {
			return nil, ErrPriceListNotFound
		}
		if e.Code == "price-lists/price-list-code-exists" {
			return nil, ErrPriceListCodeExists
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var v PriceList
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrapf(err, "decode")
	}
	return &v, nil

}

// GetPriceList calls the API service to get a price list by id.
func (c *EcomClient) GetPriceList(ctx context.Context, priceListID string) (*PriceList, error) {
	url := fmt.Sprintf("%s/price-lists/%s", c.endpoint, priceListID)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		if e.Code == "price-lists/price-list-not-found" {
			return nil, ErrPriceListNotFound
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var v PriceList
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrapf(err, "decode")
	}
	return &v, nil
}

// GetPriceLists returns a list of price lists.
func (c *EcomClient) GetPriceLists(ctx context.Context) ([]*PriceList, error) {
	url := fmt.Sprintf("%s/price-lists", c.endpoint)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	var container PriceListContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return container.Data, nil
}

// UpdatePriceList calls the API service to attempt to update a price list by id
func (c *EcomClient) UpdatePriceList(ctx context.Context, priceListID string, req *UpdatePriceListRequest) (*PriceList, error) {
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}
	body := bytes.NewReader(request)
	url := fmt.Sprintf("%s/price-lists/%s", c.endpoint, priceListID)
	res, err := c.request(http.MethodPut, url, body)
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
			return nil, errors.Wrap(err, "decode")
		}
		if e.Code == "price-lists/price-list-not-found" {
			return nil, ErrPriceListNotFound
		}
		if e.Code == "price-lists/price-list-code-exists" {
			return nil, ErrPriceListCodeExists
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var v PriceList
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrapf(err, "decode")
	}
	return &v, nil
}

// DeletePriceList calls the API service to attempt to delete a price list by id.
func (c *EcomClient) DeletePriceList(ctx context.Context, priceListID string) error {
	url := fmt.Sprintf("%s/price-lists/%s", c.endpoint, priceListID)
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrapf(err, "request(http.MethodDelete, url=%q, nil)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrap(err, "decode")
		}
		if e.Code == "price-lists/price-list-not-found" {
			return ErrPriceListNotFound
		}
		return fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}
	return nil
}
