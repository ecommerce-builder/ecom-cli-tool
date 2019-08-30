package eclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ErrProductNotFound indicates the product with given SKU could not be found.
var ErrProductNotFound = errors.New("product not found")

// ProductRequest contains fields used when applying a product.
type ProductRequest struct {
	Path string `json:"path"`
	SKU  string `json:"sku"`
	Name string `json:"name"`
}

// ProductImageApply contains the product image data.
type ProductImageApply struct {
	Path  string `json:"path" yaml:"path"`
	Title string `json:"title" yaml:"title"`
}

// ProductPricingApply contains the product pricing data.
type ProductPricingApply struct {
	TierRef   string `json:"tier_ref" yaml:"tier_ref"`
	UnitPrice int    `json:"unit_price" yaml:"unit_price"`
}

// ProductContainerResponse is a container for a list of products
type ProductContainerResponse struct {
	Object string             `json:"object"`
	Data   []*ProductResponse `json:"data"`
}

// ProductResponse contains all the fields that comprise a product in the catalog.
type ProductResponse struct {
	Object   string    `json:"object"`
	ID       string    `json:"id"`
	Path     string    `json:"path"`
	SKU      string    `json:"sku"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// ProductImage struct for capturing OpReplaceProduct JSON response.
type ProductImage struct {
	UUID     string    `json:"uuid"`
	SKU      string    `json:"sku"`
	Path     string    `json:"path"`
	GSURL    string    `json:"gsurl"`
	Width    uint      `json:"width"`
	Height   uint      `json:"height"`
	Size     uint      `json:"size"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// ProductPricing struct
type ProductPricing struct {
	UnitPrice float64   `json:"unit_price" yaml:"unit_price"`
	Created   time.Time `json:"created" yaml:"created"`
	Modified  time.Time `json:"modified" yaml:"modified"`
}

func (c *EcomClient) CreateProduct(product *ProductRequest) (*ProductResponse, error) {
	request, err := json.Marshal(&product)
	if err != nil {
		return nil, errors.Wrapf(err, "client: json marshal failed")
	}

	uri := c.endpoint + "/products"
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, uri, body)
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

	var pr ProductResponse
	if err = json.NewDecoder(res.Body).Decode(&pr); err != nil {
		return nil, errors.Wrapf(err, "replce product response decode failed")
	}
	return &pr, nil
}

// ReplaceProduct calls the API Service creating a new product or, if
// the product already exists updating it.
func (c *EcomClient) ReplaceProduct(productID string, product *ProductRequest) (*ProductResponse, error) {
	request, err := json.Marshal(&product)
	if err != nil {
		return nil, errors.Wrapf(err, "update product productID=%s failed", productID)
	}

	uri := c.endpoint + "/products/" + productID
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPut, uri, body)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, errors.Errorf("HTTP PUT to %q return %s", uri, res.Status)
	}
	var pr ProductResponse
	if err = json.NewDecoder(res.Body).Decode(&pr); err != nil {
		return nil, errors.Wrapf(err, "replce product response decode failed")
	}
	return &pr, nil
}

// GetProduct calls the API Service to get a product by SKU.
func (c *EcomClient) GetProduct(sku string) (*ProductResponse, error) {
	uri := c.endpoint + "/products/" + sku
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrProductNotFound
	} else if res.StatusCode >= 400 {
		return nil, errors.Errorf("HTTP GET to %q returned %s", uri, res.Status)
	}
	var p ProductResponse
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, errors.Wrapf(err, "get product response decode failed")
	}
	return &p, nil
}

// GetProducts returns a list of products
func (c *EcomClient) GetProducts() (*ProductContainerResponse, error) {
	uri := c.endpoint + "/products"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	var list ProductContainerResponse
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		return nil, errors.Wrapf(err, "get product response decode failed")
	}
	return &list, nil
}

func (c *EcomClient) request(method, uri string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.jwt)
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "do HTTP %s request failed", req.Method)
	}
	return res, nil
}

// ProductExists returns true if the product with the SKU sku exists.
func (c *EcomClient) ProductExists(sku string) (bool, error) {
	uri := c.endpoint + "/products/" + sku
	res, err := c.request(http.MethodHead, uri, nil)
	if err != nil {
		return false, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	if res.StatusCode == 204 {
		return true, nil
	} else if res.StatusCode == 404 {
		return false, nil
	}
	return true, errors.Wrapf(err, "unknown HTTP status code returned (%d)", res.StatusCode)
}

// DeleteProduct calls the API Service to delete a product resource.
func (c *EcomClient) DeleteProduct(sku string) error {
	exists, err := c.ProductExists(sku)
	if err != nil {
		return errors.Wrap(err, "product exists failed")
	}
	if !exists {
		return nil
	}
	uri := c.endpoint + "/products/" + sku
	res, err := c.request(http.MethodHead, uri, nil)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	return nil
}
