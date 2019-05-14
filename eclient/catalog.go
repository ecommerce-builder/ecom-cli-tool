package eclient

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// A Category is a single node in the catalog.
type Category struct {
	Segment    string     `json:"segment" yaml:"segment"`
	Name       string     `json:"name" yaml:"name"`
	Categories []Category `json:"categories" yaml:"categories"`
}

// A Catalog contains a single root node of the catalog.
type Catalog struct {
	Category Category `yaml:"catalog"`
}

// UpdateCatalog calls the API Service to update the catalog.
func (c *EcomClient) UpdateCatalog(cats Category) error {
	payload, err := json.Marshal(&cats)
	if err != nil {
		return errors.Wrapf(err, "client: json marshal failed")
	}
	uri := c.endpoint + "/catalog"
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
		return errors.Errorf("HTTP PUT to %q return %s", uri, res.Status)
	}
	return nil
}

// PurgeCatalog calls the API Service to purge the entire catalog.
func (c *EcomClient) PurgeCatalog() error {
	uri := c.endpoint + "/catalog"
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	return nil
}
