package eclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// ErrCartNotFound error
var ErrCartNotFound = errors.New("cart not found")

// ErrCartProductExists error
var ErrCartProductExists = errors.New("cart product already exists")

// ErrCartProductNotFound error
var ErrCartProductNotFound = errors.New("cart product not found")

// Cart JSON response.
type Cart struct {
	Object   string    `json:"object"`
	ID       string    `json:"id"`
	Locked   bool      `json:"locked"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// CartProduct JSON response
type CartProduct struct {
	Object    string    `json:"object"`
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Qty       int       `json:"qty"`
	UnitPrice int       `json:"unit_price"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
}

// CartProductsContainer container
type CartProductsContainer struct {
	Object string         `json:"object"`
	Data   []*CartProduct `json:"data"`
}

// CartProductRequest JSON request for adding a product to an existing cart.
type CartProductRequest struct {
	CartID    string `json:"cart_id"`
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
}

// CreateCart calls the API service to attempt to create a new shopping cart.
func (c *EcomClient) CreateCart(ctx context.Context) (*Cart, error) {
	url := c.endpoint + "/carts"
	res, err := c.request(http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: request", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	var cart Cart
	if err = json.NewDecoder(res.Body).Decode(&cart); err != nil {
		return nil, fmt.Errorf("%w: decode", err)
	}
	return &cart, nil
}

// CartAddProduct calls the API service to attempt to add an product to an
// existing cart.
func (c *EcomClient) CartAddProduct(ctx context.Context, req *CartProductRequest) (*CartProduct, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(&req); err != nil {
		return nil, fmt.Errorf("%w: encode", err)
	}
	url := c.endpoint + "/carts-products"
	res, err := c.request(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w: request", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}

		if e.Code == "carts/cart-not-found" {
			return nil, ErrCartNotFound
		}
		if e.Code == "carts/cart-product-exists" {
			return nil, ErrCartProductExists
		}
		return nil, fmt.Errorf("Status: %d, Code: %s, Message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	if res.StatusCode == 201 {
		var cartProduct CartProduct
		if err := json.NewDecoder(res.Body).Decode(&cartProduct); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		return &cartProduct, nil
	}
	return nil, fmt.Errorf("unknown response status code %d", res.StatusCode)
}

// GetCartProducts returns a list of all products in a given cart.
func (c *EcomClient) GetCartProducts(ctx context.Context, cartID string) ([]*CartProduct, error) {
	v := url.Values{}
	v.Set("cart_id", cartID)

	url := url.URL{
		Scheme:   "https",
		Host:     c.hostname,
		Path:     "carts-products",
		RawQuery: v.Encode(),
	}
	res, err := c.request(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var container CartProductsContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return container.Data, nil
}

// UpdateCartProduct calls the API service to attempt to update the qty of
// a particular product in the cart.
func (c *EcomClient) UpdateCartProduct(ctx context.Context, cartProductID string, qty int) (*CartProduct, error) {
	var req struct {
		Qty int `json:"qty"`
	}
	req.Qty = qty

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(&req); err != nil {
		return nil, fmt.Errorf("%w: encode", err)
	}
	url := fmt.Sprintf("%s/carts-products/%s", c.endpoint, cartProductID)
	res, err := c.request(http.MethodPatch, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w: request", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}

		if e.Code == "carts/cart-product-not-found" {
			return nil, ErrCartProductNotFound
		}
		return nil, fmt.Errorf("Status: %d, Code: %s, Message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	if res.StatusCode == 200 {
		var cartProduct CartProduct
		if err := json.NewDecoder(res.Body).Decode(&cartProduct); err != nil {
			return nil, errors.Wrapf(err, "json decode url %s", url)
		}
		return &cartProduct, nil
	}

	return nil, fmt.Errorf("unknown response status code %d", res.StatusCode)
}

// CartsRemoveProduct calls the API service to attempt to remove a product from a cart.
func (c *EcomClient) CartsRemoveProduct(ctx context.Context, cartProductID string) error {
	url := fmt.Sprintf("%s/carts-products/%s", c.endpoint, cartProductID)
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("%w: request", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("decode: %w", err)
		}

		if e.Code == "carts/cart-product-not-found" {
			return ErrCartProductNotFound
		}
		return fmt.Errorf("Status: %d, Code: %s, Message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	if res.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("unknown response status code %d", res.StatusCode)
}

// EmptyCartProducts calls the API service to attempt to empty all products
// from the cart
func (c *EcomClient) EmptyCartProducts(ctx context.Context, cartID string) error {
	v := url.Values{}
	v.Set("cart_id", cartID)

	url := url.URL{
		Scheme:   "https",
		Host:     c.hostname,
		Path:     "carts-products",
		RawQuery: v.Encode(),
	}
	res, err := c.request(http.MethodDelete, url.String(), nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("decode: %w", err)
		}

		if e.Code == "carts/cart-not-found" {
			return ErrCartNotFound
		}
		return fmt.Errorf("Status: %d, Code: %s, Message: %s: %w", e.Status, e.Code, e.Message, err)
	}

	if res.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("unknown response status code %d", res.StatusCode)
}
