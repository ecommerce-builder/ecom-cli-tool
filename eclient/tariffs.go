package eclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ErrShippingTariffNotFound 404
var ErrShippingTariffNotFound = errors.New("eclient: shipping tariff not found")

// ShippingTariffContainer promo rules container JSON body.
type ShippingTariffContainer struct {
	Object string            `json:"object"`
	Data   []*ShippingTariff `json:"data"`
}

// ShippingTariff shipping tariff JSON response body.
type ShippingTariff struct {
	Object       string    `json:"object"`
	ID           string    `json:"id"`
	CountryCode  string    `json:"country_code"`
	ShippingCode string    `json:"shipping_code"`
	Name         string    `json:"name"`
	Price        int       `json:"price"`
	TaxCode      string    `json:"tax_code"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
}

// CreateShippingTariffRequest request body
type CreateShippingTariffRequest struct {
	CountryCode  string `json:"country_code"`
	Shippingcode string `json:"shipping_code"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	TaxCode      string `json:"tax_code"`
}

// CreateShippingTariff call the API service to attempt to create a new
// shipping tariff.
func (c *EcomClient) CreateShippingTariff(ctx context.Context, req *CreateShippingTariffRequest) (*ShippingTariff, error) {
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("%w: client: json marshal", err)
	}

	url := c.endpoint + "/shipping-tariffs"
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

	var tariff ShippingTariff
	if err = json.NewDecoder(res.Body).Decode(&tariff); err != nil {
		return nil, fmt.Errorf("%w: decode", err)
	}
	return &tariff, nil
}

// GetShippingTariffs returns a list of all shipping tariffs.
func (c *EcomClient) GetShippingTariffs(ctx context.Context) ([]*ShippingTariff, error) {
	uri := c.endpoint + "/shipping-tariffs"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var container ShippingTariffContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("get price list response decode failed: %w", err)
	}
	return container.Data, nil
}
