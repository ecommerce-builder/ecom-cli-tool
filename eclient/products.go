package eclient

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type ProductCreate struct {
	SKU  string      `json:"sku" yaml:"sku"`
	EAN  string      `json:"ean" yaml:"ean"`
	URL  string      `json:"url" yaml:"url"`
	Name string      `json:"name" yaml:"name"`
	Data ProductData `json:"data" yaml:"data"`
}

type ProductData struct {
	Summary string `json:"summary" yaml:"summary"`
	Desc    string `json:"description" yaml:"description"`
	Spec    string `json:"specification" yaml:"specification"`
}

type ProductUpdate struct {
	EAN  string      `json:"ean" yaml:"ean"`
	URL  string      `json:"url" yaml:"url"`
	Name string      `json:"name" yaml:"name"`
	Data ProductData `json:"data" yaml:"data"`
}

// Product contains all the fields that comprise a product in the catalog.
type Product struct {
	SKU      string      `json:"sku" yaml:"sku,omitempty"`
	EAN      string      `json:"ean" yaml:"ean"`
	URL      string      `json:"url" yaml:"url"`
	Name     string      `json:"name" yaml:"name"`
	Data     ProductData `json:"data" yaml:"data"`
	Created  time.Time   `json:"created,omitempty"`
	Modified time.Time   `json:"modified,omitempty"`
}

// UpdateProduct calls the API Service to update an existing product.
func (c *EcomClient) UpdateProduct(sku string, p *ProductUpdate) (*Product, error) {
	payload, err := json.Marshal(&p)
	if err != nil {
		return nil, errors.Wrapf(err, "update product sku=%q failed", sku)
	}
	uri := c.endpoint + "/products/" + sku
	body := strings.NewReader(string(payload))
	res, err := c.request(http.MethodPut, uri, body)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	//bs, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return nil, errors.Wrap(err, "readall failed:")
	//}
	//fmt.Println(string(bs))
	if res.StatusCode >= 400 {
		return nil, errors.Errorf("HTTP PUT to %q return %s", uri, res.Status)
	}
	pr := Product{}
	err = json.NewDecoder(res.Body).Decode(&pr)
	if err != nil {
		return nil, errors.Wrapf(err, "create product response decode failed")
	}
	return &pr, nil
}

// CreateProduct calls the API Service to create a new product.
func (c *EcomClient) CreateProduct(p *ProductCreate) (*Product, error) {
	payload, err := json.Marshal(&p)
	if err != nil {
		return nil, errors.Wrapf(err, "create product sku=%q failed", p.SKU)
	}
	uri := c.endpoint + "/products"
	body := strings.NewReader(string(payload))
	res, err := c.request(http.MethodPost, uri, body)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, errors.Errorf("HTTP POST to %q return %s", uri, res.Status)
	}
	productResponse := Product{}
	err = json.NewDecoder(res.Body).Decode(&productResponse)
	if err != nil {
		return nil, errors.Wrapf(err, "create product response decode failed")
	}
	return &productResponse, nil
}

// GetProduct calls the API Service to get a product by SKU.
func (c *EcomClient) GetProduct(sku string) (*Product, error) {
	uri := c.endpoint + "/products/" + sku
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	p := Product{}
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
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
