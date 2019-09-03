package eclient

import (
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

// PricesContainerResponse a container for a list of price lists.
type PricesContainerResponse struct {
	Object string           `json:"object"`
	Data   []*PriceResponse `json:"data"`
}

// PriceResponse a single price.
type PriceResponse struct {
	Object      string    `json:"object"`
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	PriceListID string    `json:"price_list_id"`
	Break       int       `json:"break"`
	UnitPrice   int       `json:"unit_price"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// SetPrices calls the API Service to update the categories tree.
func (c *EcomClient) SetPrices(productID, priceListID string, prices []*PriceRequest) ([]*PriceResponse, error) {
	container := PricesContainerRequest{
		Object: "list",
		Data:   prices,
	}

	request, err := json.Marshal(&container)
	if err != nil {
		return nil, errors.Wrapf(err, "client: json marshal failed")
	}

	params := url.Values{}
	params.Add("product_id", productID)
	params.Add("price_list_id", priceListID)

	uri := c.endpoint + "/prices?" + params.Encode()
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPut, uri, body)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, errors.Wrapf(err, "client decode error")
		}
		return nil, errors.Errorf(fmt.Sprintf("Status: %d, Code: %s, Message: %s\n", e.Status, e.Code, e.Message))
	}

	var response PricesContainerResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrapf(err, "json decode url %s failed", uri)
	}
	return response.Data, nil
}
