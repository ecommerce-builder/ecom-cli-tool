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

// PricesContainerRequest JSON request body.
type PricesContainerRequest struct {
	Object string          `json:"object"`
	Data   []*PriceRequest `json:"data"`
}

// PriceRequest JSON price request body.
type PriceRequest struct {
	Break     int `json:"break"`
	UnitPrice int `json:"unit_price"`
}

// PricesContainer a container for a list of price lists.
type PricesContainer struct {
	Object string   `json:"object"`
	Data   []*Price `json:"data"`
}

// Price a single price.
type Price struct {
	Object        string    `json:"object"`
	ID            string    `json:"id"`
	ProductID     string    `json:"product_id"`
	ProductPath   string    `json:"product_path"`
	ProductSKU    string    `json:"product_sku"`
	PriceListID   string    `json:"price_list_id"`
	PriceListCode string    `json:"price_list_code"`
	Break         int       `json:"break"`
	UnitPrice     int       `json:"unit_price"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modified"`
}

// SetPrices calls the API Service to update the categories tree.
func (c *EcomClient) SetPrices(productID, priceListID string, prices []*PriceRequest) ([]*Price, error) {
	container := PricesContainerRequest{
		Object: "list",
		Data:   prices,
	}

	request, err := json.Marshal(&container)
	if err != nil {
		return nil, fmt.Errorf("client: json marshal failed: %w", err)
	}

	params := url.Values{}
	params.Add("product_id", productID)
	params.Add("price_list_id", priceListID)

	uri := c.endpoint + "/prices?" + params.Encode()
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPut, uri, body)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("client decode error: %w", err)
		}
		return nil, fmt.Errorf("Status: %d, Code: %s, Message: %s", e.Status, e.Code, e.Message)
	}

	var response PricesContainer
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrapf(err, "json decode url %s", uri)
	}
	return response.Data, nil
}

// GetPrices calls the API service to attempt to retrieve all prices
// for all products.
func (c *EcomClient) GetPrices(ctx context.Context) ([]*Price, error) {
	url := fmt.Sprintf("%s/prices", c.endpoint)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err,
			"request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	var container PricesContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return container.Data, nil
}
