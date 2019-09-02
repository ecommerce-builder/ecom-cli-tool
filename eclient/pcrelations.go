package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ProductCategoryContainerResponse container.
type ProductCategoryContainerResponse struct {
	Object string                     `json:"object"`
	Data   []*ProductCategoryResponse `json:"data"`
}

// ProductCategoryResponse response body.
type ProductCategoryResponse struct {
	Object       string    `json:"object"`
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	ProductPath  string    `json:"product_path"`
	ProductSKU   string    `json:"product_sku"`
	ProductName  string    `json:"product_name"`
	CategoryID   string    `json:"category_id"`
	CategoryPath string    `json:"category_path"`
	Pri          int       `json:"pri"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
}

// CreateProductsCategoriesContainer request body
type CreateProductsCategoriesContainer struct {
	Object string                      `json:"object"`
	Data   []*CreateProductsCategories `json:"data"`
}

// CreateProductsCategories request body
type CreateProductsCategories struct {
	ProductID  string `json:"product_id"`
	CategoryID string `json:"category_id"`
}

// An AssocProduct holds details of a product in the context of an AssocSet.
type AssocProduct struct {
	SKU      string    `json:"sku" yaml:"sku"`
	Created  time.Time `json:"created,omitempty"`
	Modified time.Time `json:"modified,omitempty"`
}

// AssocResponse details a catalog association including products.
type AssocResponse struct {
	Products []AssocProduct `json:"products"`
}

// GetProductCategoryRelations calls the API Service to get all catalog associations.
func (c *EcomClient) GetProductCategoryRelations() ([]*ProductCategoryResponse, error) {
	uri := c.endpoint + "/products-categories"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	var container ProductCategoryContainerResponse
	err = json.NewDecoder(res.Body).Decode(&container)
	if err != nil {
		return nil, errors.Wrapf(err, "get product response decode failed")
	}
	return container.Data, nil
}

// UpdateProductCategoryRelations calls the API Service to update all
// product to category relations.
func (c *EcomClient) UpdateProductCategoryRelations(rels []*CreateProductsCategories) error {
	container := CreateProductsCategoriesContainer{
		Object: "list",
		Data:   rels,
	}

	payload, err := json.Marshal(&container)
	if err != nil {
		return errors.Wrapf(err, "client: json marshal failed")
	}
	uri := c.endpoint + "/products-categories"
	body := strings.NewReader(string(payload))
	res, err := c.request(http.MethodPut, uri, body)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	//bs, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return nil, errors.Wrap(err, "readall failed:")
	//}
	//fmt.Println(string(bs))
	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, "client decode error")
		}
		return errors.Errorf(fmt.Sprintf("Status: %d, Code: %s, Message: %s\n", e.Status, e.Code, e.Message))
	}
	return nil
}

// DeleteProductCategoryRelations calls the API Service to delete all
// product to category relations.
func (c *EcomClient) DeleteProductCategoryRelations() error {
	uri := c.endpoint + "/products-categories"
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return errors.Errorf("unauthorized")
	}
	return nil
}
