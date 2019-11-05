package eclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ErrCouponNotFound coupon not found
var ErrCouponNotFound = errors.New("coupon not found")

// Coupon a single coupon for use with the cart.
type Coupon struct {
	Object        string    `json:"object"`
	ID            string    `json:"id"`
	CouponCode    string    `json:"coupon_code"`
	PromoRuleID   string    `json:"promo_rule_id"`
	PromoRuleCode string    `json:"promo_rule_code"`
	Void          bool      `json:"void"`
	Resuable      bool      `json:"resuable"`
	SpendCount    int       `json:"spend_count"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modfied"`
}

type couponsListContainer struct {
	Object string    `json:"object"`
	Data   []*Coupon `json:"data"`
}

// CreateCouponRequest request body
type CreateCouponRequest struct {
	PromoRuleID string `json:"promo_rule_id"`
	CouponCode  string `json:"coupon_code"`
	Resuable    bool   `json:"reusable"`
}

// CreateCoupon calls the API service to attempt to mint a new coupon.
func (c *EcomClient) CreateCoupon(ctx context.Context, req *CreateCouponRequest) (*Coupon, error) {
	request, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrapf(err, "json marshal")
	}

	url := fmt.Sprintf("%s/coupons", c.endpoint)
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "request")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, errors.Wrap(err, "decode")
		}
		return nil, errors.Wrapf(err,
			"status: %d, code: %s, message: %s: %v",
			e.Status, e.Code, e.Message, err)
	}

	var coupon Coupon
	if err = json.NewDecoder(res.Body).Decode(&coupon); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &coupon, nil
}

// GetCoupons calls the API service to attempt to retrieve all coupons.
func (c *EcomClient) GetCoupons(ctx context.Context) ([]*Coupon, error) {
	url := fmt.Sprintf("%s/coupons", c.endpoint)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "request(http.MethodGet, url=%q, nil)", url)
	}
	defer res.Body.Close()

	var container couponsListContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}
	return container.Data, nil
}

// DeleteCoupon calls the API service to attempt to delete a coupon
// with the given id.
func (c *EcomClient) DeleteCoupon(ctx context.Context, couponID string) error {
	url := fmt.Sprintf("%s/coupons/%s", c.endpoint, couponID)
	res, err := c.request(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrapf(err, "delete request url=%q", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return errors.Wrap(err, "decode")
		}
		if e.Code == "coupons/coupon-not-found" {
			return ErrCouponNotFound
		}
		return errors.Wrapf(err,
			"status: %d, code: %s, message: %s: %v",
			e.Status, e.Code, e.Message, err)
	}
	return nil
}

// VoidCoupon calls the API service to attempt to void a coupon.
func (c *EcomClient) VoidCoupon(ctx context.Context, couponID string) error {
	request, err := json.Marshal(struct {
		Void bool `json:"void"`
	}{
		Void: true,
	})
	if err != nil {
		return errors.Wrapf(err, "json marshal")
	}

	body := strings.NewReader(string(request))
	url := fmt.Sprintf("%s/coupons/%s", c.endpoint, couponID)
	res, err := c.request(http.MethodPatch, url, body)
	if err != nil {
		return errors.Wrapf(err, "patch request url=%q", url)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return errors.Wrap(err, "decode")
		}
		if e.Code == "coupons/coupon-not-found" {
			return ErrCouponNotFound
		}
		return errors.Wrapf(err,
			"status: %d, code: %s, message: %s: %v",
			e.Status, e.Code, e.Message, err)
	}
	return nil
}
