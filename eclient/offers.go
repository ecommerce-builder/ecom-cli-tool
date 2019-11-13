package eclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// ErrOfferNotFound error (offers/offer-not-found)
var ErrOfferNotFound = errors.New("offer not found")

// ErrOfferExists error (offers/offer-exists)
var ErrOfferExists = errors.New("offer exists")

// Offer response
type Offer struct {
	Object        string    `json:"object"`
	ID            string    `json:"id"`
	PromoRuleID   string    `json:"promo_rule_id"`
	PromoRuleCode string    `json:"promo_rule_code"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modfied"`
}

// OfferContainer list container
type OfferContainer struct {
	Object string   `json:"object"`
	Data   []*Offer `json:"data"`
}

// CreateOfferRequest request
type CreateOfferRequest struct {
	PromoRuleID string `json:"promo_rule_id"`
}

// CreateOffer calls the API service to active an offer using the given
// promo rule id.
func (c *EcomClient) CreateOffer(ctx context.Context, req *CreateOfferRequest) (*Offer, error) {
	fmt.Printf("%+v\n", req.PromoRuleID)
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrapf(err, "json marshal")
	}
	body := bytes.NewReader(request)
	url := fmt.Sprintf("%s/offers", c.endpoint)
	res, err := c.request(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrapf(err,
			"request(http.MethodPost, url=%q, body=%q)",
			url, string(request))
	}
	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		if e.Code == "offers/offer-exists" {
			return nil, ErrOfferExists
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}

	var v Offer
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &v, nil
}

// GetOffer calls the API service to get an individual offer.
func (c *EcomClient) GetOffer(ctx context.Context, offerID string) (*Offer, error) {
	url := fmt.Sprintf("%s/offers/%s", c.endpoint, offerID)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		if e.Code == "offers/offer-not-found" {
			return nil, ErrOfferNotFound
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}

	var v Offer
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &v, nil
}

// GetOffers calls the API service to get a list of all active offers.
func (c *EcomClient) GetOffers(ctx context.Context) ([]*Offer, error) {
	url := fmt.Sprintf("%s/offers", c.endpoint)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}

	var container OfferContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return container.Data, nil
}

// DeleteOffer calls the API service to delete and offer by id.
func (c *EcomClient) DeleteOffer(ctx context.Context, offerID string) error {
	url := fmt.Sprintf("%s/webhooks/%s", c.endpoint, offerID)
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrapf(err, "request(http.MethodDelete, url=%q, nil)", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return errors.Wrap(err, "decode")
		}
		if e.Code == "offers/offer-not-found" {
			return ErrOfferNotFound
		}
		return fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}
	return nil
}
