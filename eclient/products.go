package eclient

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ErrProductNotFound indicates the product with given SKU could not be found.
var ErrProductNotFound = errors.New("product not found")

// ProductContent contains the variable JSON data of the product
type ProductContent struct {
	Meta struct {
		Title       string `json:"title" yaml:"title"`
		Description string `json:"description" yaml:"description"`
	} `json:"meta" yaml:"meta"`
	Videos        []string `json:"videos" yaml:"videos"`
	Manuals       []string `json:"manuals" yaml:"manuals"`
	Software      []string `json:"software" yaml:"software"`
	Description   string   `json:"description" yaml:"description"`
	Specification string   `json:"specification" yaml:"specification"`
	InTheBox      string   `json:"in_the_box" yaml:"in_the_box"`
}

// ProductApply contains fields used when applying a product.
type ProductApply struct {
	SKU     string                 `json:"sku,omitempty" yaml:"sku"`
	EAN     string                 `json:"ean" yaml:"ean"`
	Path    string                 `json:"path" yaml:"path"`
	Name    string                 `json:"name" yaml:"name"`
	Images  []*ProductImageApply   `json:"images" yaml:"images"`
	Pricing []*ProductPricingApply `json:"pricing" yaml:"pricing"`
	Content ProductContent         `json:"content" yaml:"content"`
}

// ProductImageApply contains the product image data.
type ProductImageApply struct {
	Path  string `json:"path" yaml:"path"`
	Title string `json:"title" yaml:"title"`
}

// ProductPricingApply contains the product pricing data.
type ProductPricingApply struct {
	TierRef   string  `json:"tier_ref" yaml:"tier_ref"`
	UnitPrice float64 `json:"unit_price" yaml:"unit_price"`
}

// Product contains all the fields that comprise a product in the catalog.
type Product struct {
	SKU      string                     `json:"sku" yaml:"sku"`
	EAN      string                     `json:"ean" yaml:"ean"`
	Path     string                     `json:"path" yaml:"path"`
	Name     string                     `json:"name" yaml:"name"`
	Images   []*ProductImage            `json:"images" yaml:"images"`
	Pricing  map[string]*ProductPricing `json:"pricing" yaml:"pricing"`
	Content  ProductContent             `json:"content" yaml:"content"`
	Created  time.Time                  `json:"created,omitempty"`
	Modified time.Time                  `json:"modified,omitempty"`
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

// ProductPricing
type ProductPricing struct {
	UnitPrice float64   `json:"unit_price" yaml:"unit_price"`
	Created   time.Time `json:"created" yaml:"created"`
	Modified  time.Time `json:"modified" yaml:"modified"`
}

// ReplaceProduct calls the API Service creating a new product or, if
// the product already exists updating it.
func (c *EcomClient) ReplaceProduct(sku string, p *ProductApply) (*Product, error) {
	payload, err := json.Marshal(&p)
	if err != nil {
		return nil, errors.Wrapf(err, "update product sku=%s failed", sku)
	}
	uri := c.endpoint + "/products/" + sku
	body := strings.NewReader(string(payload))
	res, err := c.request(http.MethodPut, uri, body)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, errors.Errorf("HTTP PUT to %q return %s", uri, res.Status)
	}
	pr := Product{}
	if err = json.NewDecoder(res.Body).Decode(&pr); err != nil {
		return nil, errors.Wrapf(err, "replce product response decode failed")
	}
	return &pr, nil
}

// GetProduct calls the API Service to get a product by SKU.
func (c *EcomClient) GetProduct(sku string) (*Product, error) {
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
	p := Product{}
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, errors.Wrapf(err, "get product response decode failed")
	}
	return &p, nil
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
