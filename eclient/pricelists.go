package eclient

import (
	"encoding/json"
	"time"

	"net/http"

	"github.com/pkg/errors"
)

// PriceListContainerResponse price list container JSON body.
type PriceListContainerResponse struct {
	Object string               `json:"object"`
	Data   []*PriceListResponse `json:"data"`
}

// PriceListResponse price list JSON response body.
type PriceListResponse struct {
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

// GetPriceLists returns a list of price lists.
func (c *EcomClient) GetPriceLists() ([]*PriceListResponse, error) {
	uri := c.endpoint + "/price-lists"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	var container PriceListContainerResponse
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, errors.Wrapf(err, "get price list response decode failed")
	}
	return container.Data, nil
}