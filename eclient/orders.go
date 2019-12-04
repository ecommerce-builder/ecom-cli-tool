package eclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ErrOrderNotFound error
var ErrOrderNotFound = errors.New("order not found")

// Order summary
type Order struct {
	Object   string       `json:"object"`
	ID       string       `json:"id"`
	OrderID  int          `json:"order_id"`
	Status   string       `json:"status"`
	Payment  string       `json:"payment"`
	Billing  OrderAddr    `json:"billing_address"`
	Shipping OrderAddr    `json:"shipping_address"`
	Items    []*OrderItem `json:"items"`
	User     struct {
		ContactName string `json:"contact_name,omitempty"`
		Email       string `json:"email,omitempty"`
		UserID      string `json:"user_id,omitempty"`
	} `json:"user"`
	Currency    string    `json:"currency"`
	TotalExVAT  int       `json:"total_ex_vat"`
	VATTotal    int       `json:"vat_total"`
	TotalIncVAT int       `json:"total_inc_vat"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// OrderAddr shipping or billing address
type OrderAddr struct {
	ContactName string  `json:"contact_name"`
	Addr1       string  `json:"addr1"`
	Addr2       *string `json:"addr2"`
	City        string  `json:"city"`
	County      *string `json:"county"`
	Postcode    string  `json:"postcode"`
	CountryCode string  `json:"country_code"`
}

// OrderItem product item inside the order
type OrderItem struct {
	Object    string    `json:"object"`
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Qty       int       `json:"qty"`
	UnitPrice int       `json:"unit_price"`
	Currency  string    `json:"currency"`
	TaxCode   string    `json:"tax_code"`
	VAT       int       `json:"vat"`
	Created   time.Time `json:"created"`
}

// OrderListContainer object
type OrderListContainer struct {
	Object string   `json:"object"`
	Data   []*Order `json:"data"`
}

// OrderRequest for guest and registered user orders.
type OrderRequest struct {
	CartID      *string              `json:"cart_id"`
	ContactName *string              `json:"contact_name,omitempty"`
	Email       *string              `json:"email,omitempty"`
	UserID      *string              `json:"user_id,omitempty"`
	BillingID   *string              `json:"billing_id,omitempty"`
	ShippingID  *string              `json:"shipping_id,omitempty"`
	Billing     *OrderAddressRequest `json:"billing,omitempty"`
	Shipping    *OrderAddressRequest `json:"shipping,omitempty"`
}

// OrderAddressRequest contains the new address request body
type OrderAddressRequest struct {
	ContactName *string `json:"contact_name"`
	Addr1       *string `json:"addr1"`
	Addr2       *string `json:"addr2"`
	City        *string `json:"city"`
	County      *string `json:"county"`
	Postcode    *string `json:"postcode"`
	CountryCode *string `json:"country_code"`
}

// PlaceOrder calls the API service to attempt to place an order
func (c *EcomClient) PlaceOrder(ctx context.Context, o *OrderRequest) (*Order, error) {
	fmt.Printf("%+v\n", o)
	request, err := json.Marshal(&o)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}

	body := strings.NewReader(string(request))
	url := fmt.Sprintf("%s/orders", c.endpoint)
	res, err := c.request(http.MethodPost, url, body)
	if err != nil {
		return nil, ErrCartNotFound
	}

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}

	var v Order
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrapf(err, "decode")
	}
	return &v, nil
}

// GetOrder calls the API service to return a single order by id.
func (c *EcomClient) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	url := fmt.Sprintf("%s/orders/%s", c.endpoint, orderID)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}

	var v Order
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &v, nil

}

// GetOrders calls the API service to return all orders.
func (c *EcomClient) GetOrders(ctx context.Context) ([]*Order, error) {
	url := fmt.Sprintf("%s/orders", c.endpoint)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err,
			"request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		//dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}

	var container OrderListContainer
	if err = json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return container.Data, nil
}
