package eclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ErrPPAssocGroupNotFound error.
var ErrPPAssocGroupNotFound = errors.New("product to product associations group not found")

// CreatePAGroupRequest Product Associations Group request
type CreatePAGroupRequest struct {
	PPAssocGroupCode string `json:"pp_assoc_group_code"`
	Name             string `json:"name"`
}

// PPAssocGroupContainer container
type PPAssocGroupContainer struct {
	Object string                  `json:"object"`
	Data   []*PPAssocGroupResponse `json:"data"`
}

// PPAssocGroupResponse holds a product to product associations group.
type PPAssocGroupResponse struct {
	Object   string    `json:"object"`
	ID       string    `json:"id"`
	Code     string    `json:"pp_assoc_group_code"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// CreatePPAGroup create a new product to product association group.
func (c *EcomClient) CreatePPAGroup(ctx context.Context, g *CreatePAGroupRequest) (*PPAssocGroupResponse, error) {
	request, err := json.Marshal(&g)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	uri := c.endpoint + "/products-assocs-groups"
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, uri, body)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrPPAssocGroupNotFound
	}

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("client decode: %w", err)
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var ppaGroup PPAssocGroupResponse
	if err = json.NewDecoder(res.Body).Decode(&ppaGroup); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &ppaGroup, nil
}

// GetPPAGroup calls the API service to get a single product to product
// associations group.
func (c *EcomClient) GetPPAGroup(ctx context.Context, ppaGroupID string) (*PPAssocGroupResponse, error) {
	url := c.endpoint + "/products-assocs-groups/" + ppaGroupID
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		return nil, ErrBadRequest
	}

	if res.StatusCode == 404 {
		return nil, ErrPPAssocGroupNotFound
	}

	var g PPAssocGroupResponse
	if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &g, nil
}

// GetPPAGroups returns a list of all product to product association groups.
func (c *EcomClient) GetPPAGroups(ctx context.Context) ([]*PPAssocGroupResponse, error) {
	uri := c.endpoint + "/products-assocs-groups"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	var container PPAssocGroupContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return container.Data, nil
}

// DeletePPAGroup calls the API service to attempt to delete a product to
// product associations group with the given ppaGroupID.
func (c *EcomClient) DeletePPAGroup(ctx context.Context, ppaGroupID string) error {
	url := c.endpoint + "/products-assocs-groups/" + ppaGroupID
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		return nil
	}
	if res.StatusCode == 400 {
		return ErrBadRequest
	}
	if res.StatusCode == 404 {
		return ErrPPAssocGroupNotFound
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("err response %d", res.StatusCode)
	}
	return nil
}
