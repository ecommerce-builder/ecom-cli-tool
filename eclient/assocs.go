package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// An Assoc holds a single catalog association.
type Assoc struct {
	Products []string `json:"products" yaml:"products"`
}

// Associations for the catalog associations.
type Associations struct {
	Assocs map[string]*Assoc `yaml:"associations"`
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

// GetCatalogAssocs calls the API Service to get all catalog associations.
func (c *EcomClient) GetCatalogAssocs() (map[string]AssocResponse, error) {
	uri := c.endpoint + "/associations"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	assocs := make(map[string]AssocResponse)
	err = json.NewDecoder(res.Body).Decode(&assocs)
	if err != nil {
		return nil, errors.Wrapf(err, "get product response decode failed")
	}
	return assocs, nil
}

// UpdateCatalogAssocs calls the API Service to update all catalog associations.
func (c *EcomClient) UpdateCatalogAssocs(assocs map[string]*Assoc) error {
	payload, err := json.Marshal(&assocs)
	if err != nil {
		return errors.Wrapf(err, "client: json marshal failed")
	}
	uri := c.endpoint + "/associations"
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
		var e struct {
			Status  int    `json:"status"`
			Code    string `json:"code"`
			Message string `json:"message"`
			Data    struct {
				MissingPaths []string `json:"missing_path"`
				NonLeafPaths []string `json:"non_leaf_paths"`
				MissingSKUs  []string `json:"missing_skus"`
			} `json:"data"`
		}
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return errors.Wrapf(err, "4xx decode error")
		}
		// fmt.Printf("%+v\n", e)
		fmt.Printf("%s\n", e.Message)
		return errors.Errorf("HTTP PUT to %q return %s", uri, res.Status)
	}
	return nil
}

// PurgeCatalogAssocs calls the API Service to delete all catalog associations.
func (c *EcomClient) PurgeCatalogAssocs() error {
	uri := c.endpoint + "/associations"
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
