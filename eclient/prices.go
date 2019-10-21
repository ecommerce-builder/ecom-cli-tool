package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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

// PricesContainerResponse a container for a list of price lists.
type PricesContainerResponse struct {
	Object string           `json:"object"`
	Data   []*PriceResponse `json:"data"`
}

// PriceResponse a single price.
type PriceResponse struct {
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
func (c *EcomClient) SetPrices(productID, priceListID string, prices []*PriceRequest) ([]*PriceResponse, error) {
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
		return nil, fmt.Errorf("Status: %d, Code: %s, Message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	var response PricesContainerResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("json decode url %s failed: %w", uri, err)
	}
	return response.Data, nil
}
