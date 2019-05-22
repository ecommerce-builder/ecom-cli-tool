package eclient

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type createAdminRequest struct {
	Email  string `json:"email"`
	Passwd string `json:"password"`
	First  string `json:"firstname"`
	Last   string `json:"lastname"`
}

// CreateAdmin calls the API Service to create a new administrator.
func (c *EcomClient) CreateAdmin(email, passwd, first, last string) (*Customer, error) {
	p := createAdminRequest{
		Email:  email,
		Passwd: passwd,
		First:  first,
		Last:   last,
	}
	payload, err := json.Marshal(&p)
	if err != nil {
		return nil, errors.Wrap(err, "create admin failed")
	}
	uri := c.endpoint + "/admins"
	res, err := c.request(http.MethodPost, uri, strings.NewReader(string(payload)))
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, errors.Errorf("HTTP POST to %q return %s", uri, res.Status)
	}
	customer := Customer{}
	err = json.NewDecoder(res.Body).Decode(&customer)
	if err != nil {
		return nil, errors.Wrapf(err, "create product response decode failed")
	}
	return &customer, nil
}

// ListAdmins calls the API Service to get all administrators.
func (c *EcomClient) ListAdmins() ([]*Customer, error) {
	uri := c.endpoint + "/admins"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	admins := make([]*Customer, 0, 8)
	err = json.NewDecoder(res.Body).Decode(&admins)
	if err != nil {
		return nil, errors.Wrapf(err, "list admins response decode failed")
	}
	return admins, nil
}

// DeleteAdmin calls the API service to delete an administrator with the
// given UUID.
func (c *EcomClient) DeleteAdmin(uuid string) error {
	uri := c.endpoint + "/admins/" + uuid
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()
	return nil
}
