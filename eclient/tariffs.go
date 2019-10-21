package eclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ErrShippingTariffNotFound 404
var ErrShippingTariffNotFound = errors.New("eclient: shipping tariff not found")

// ShippingTariffContainerResponse promo rules container JSON body.
type ShippingTariffContainerResponse struct {
	Object string                    `json:"object"`
	Data   []*ShippingTariffResponse `json:"data"`
}

// ShippingTariffResponse shipping tariff JSON response body.
type ShippingTariffResponse struct {
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

// GetShippingTariffs returns a list of all shipping tariffs.
func (c *EcomClient) GetShippingTariffs(ctx context.Context) ([]*ShippingTariffResponse, error) {
	uri := c.endpoint + "/shipping-tariffs"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var container ShippingTariffContainerResponse
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("get price list response decode failed: %w", err)
	}
	return container.Data, nil
}
