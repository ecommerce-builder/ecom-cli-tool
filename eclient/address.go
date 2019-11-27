package eclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ErrAddressNotFound error
var ErrAddressNotFound = errors.New("address not found")

// AddressListContainer object
type AddressListContainer struct {
	Object string     `json:"object"`
	Data   []*Address `json:"data"`
}

// Address contains address information for a user
type Address struct {
	Object      string    `json:"object"`
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Typ         string    `json:"type"`
	ContactName string    `json:"contact_name"`
	Addr1       string    `json:"addr1"`
	Addr2       *string   `json:"addr2,omitempty"`
	City        string    `json:"city"`
	County      *string   `json:"county,omitempty"`
	Postcode    string    `json:"postcode"`
	CountryCode string    `json:"country_code"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// CreateAddressRequest JSON payload
type CreateAddressRequest struct {
	UserID      string  `json:"user_id"`
	Type        string  `json:"type"`
	ContactName string  `json:"contact_name"`
	Addr1       string  `json:"addr1"`
	Addr2       *string `json:"addr2"`
	City        string  `json:"city"`
	County      *string `json:"county"`
	Postcode    string  `json:"postcode"`
	CountryCode string  `json:"country_code"`
}

// UpdateAddressRequest JSON payload
type UpdateAddressRequest struct {
	Type        *string `json:"type,omitempty"`
	ContactName *string `json:"contact_name,omitempty"`
	Addr1       *string `json:"addr1,omitempty"`
	Addr2       *string `json:"addr2,omitempty"`
	City        *string `json:"city,omitempty"`
	County      *string `json:"county,omitempty"`
	Postcode    *string `json:"postcode,omitempty"`
	CountryCode *string `json:"country_code,omitempty"`
}

// CreateAddress calls the API service to attempt to create a
// new address object.
func (c *EcomClient) CreateAddress(ctx context.Context, req *CreateAddressRequest) (*Address, error) {
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("%w: client: json marshal", err)
	}

	url := fmt.Sprintf("%s/addresses", c.endpoint)
	body := strings.NewReader(string(request))

	res, err := c.request(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w: request", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("%w: client decode", err)
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var address Address
	if err = json.NewDecoder(res.Body).Decode(&address); err != nil {
		return nil, fmt.Errorf("%w: decode", err)
	}
	return &address, nil
}

// GetAddress calls the API service to attempt to retrieve a single address
// by id.
func (c *EcomClient) GetAddress(ctx context.Context, addrID string) (*Address, error) {
	url := fmt.Sprintf("%s/addresses/%s", c.endpoint, addrID)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		if e.Code == "addresses/address-not-found" {
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var v Address
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &v, nil
}

// GetAddressesByUser calls the API service to attempt to retrieve all the
// addresses for a particular user.
func (c *EcomClient) GetAddressesByUser(ctx context.Context, userID string) ([]*Address, error) {
	v := url.Values{}
	v.Set("user_id", userID)
	url := url.URL{
		Scheme:   c.scheme,
		Host:     c.hostname,
		Path:     "addresses",
		RawQuery: v.Encode(),
	}
	res, err := c.request(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	var container AddressListContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return container.Data, nil
}

// UpdateAddress calls the API service to attempt to update the address
// with the given addrID.
func (c *EcomClient) UpdateAddress(ctx context.Context, addrID string, req *UpdateAddressRequest) (*Address, error) {
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrapf(err, "json marshal: %v", req)
	}

	body := strings.NewReader(string(request))
	url := c.endpoint + "/addresses/" + addrID
	res, err := c.request(http.MethodPatch, url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "c.request(http.MethodPatch, uri=%q, body)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		if e.Code == "addresses/address-not-found" {
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("Status: %d, Code: %s, Message: %s", e.Status, e.Code, e.Message)
	}

	var v Address
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrapf(err, "json decode url %s", url)
	}
	return &v, nil
}

// DeleteAddress calls the API service to attempt to delete an address.
func (c *EcomClient) DeleteAddress(ctx context.Context, addrID string) error {
	url := fmt.Sprintf("%s/addresses/%s", c.endpoint, addrID)
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		if e.Code == "addresses/address-not-found" {
			return ErrAddressNotFound
		}
		return fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}
	return nil
}
