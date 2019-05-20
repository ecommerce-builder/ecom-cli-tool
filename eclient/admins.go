package eclient

import (
	"encoding/json"
	"net/http"

	service "bitbucket.org/andyfusniakteam/ecom-api-go/service/firebase"
	"github.com/pkg/errors"
)

// ListAdmins calls the API Service to get all administrators.
func (c *EcomClient) ListAdmins() ([]*service.Customer, error) {
	uri := c.endpoint + "/admins"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	admins := make([]*service.Customer, 0, 8)
	err = json.NewDecoder(res.Body).Decode(&admins)
	if err != nil {
		return nil, errors.Wrapf(err, "list admins response decode failed")
	}
	return admins, nil
}
