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

// ErrBadRequest is returned on 400 status code
var ErrBadRequest = errors.New("eclient: bad request")

// ErrPromoRuleNotFound 404
var ErrPromoRuleNotFound = errors.New("eclient: promo rule not found")

// PromoRulesContainer promo rules container JSON body.
type PromoRulesContainer struct {
	Object string       `json:"object"`
	Data   []*PromoRule `json:"data"`
}

// PromoRuleRequest request body
type PromoRuleRequest struct {
	PromoRuleCode    string     `json:"promo_rule_code"`
	Name             string     `json:"name"`
	StartAt          *time.Time `json:"start_at,omitempty"`
	EndAt            *time.Time `json:"end_at,omitempty"`
	Amount           int        `json:"amount,omitempty"`
	TotalThreshold   int        `json:"total_threshold,omitempty"`
	Type             string     `json:"type"`
	Target           string     `json:"target"`
	ProductID        string     `json:"product_id,omitempty"`
	CategoryID       string     `json:"category_id,omitempty"`
	ShippingTariffID string     `json:"shipping_tariff_id,omitempty"`
}

// PromoRule price list JSON response body.
type PromoRule struct {
	Object             string     `json:"object"`
	ID                 string     `json:"id"`
	PromoRuleCode      string     `json:"promo_rule_code"`
	ProductID          *string    `json:"product_id,omitempty"`
	ProductPath        *string    `json:"product_path,omitempty"`
	ProductSKU         *string    `json:"product_sku,omitempty"`
	CategoryID         *string    `json:"category_id,omitempty"`
	CategoryPath       *string    `json:"category_path,omitempty"`
	ShippingTariffID   *string    `json:"shipping_tariff_id,omitempty"`
	ShippingTariffCode *string    `json:"shipping_tariff_code,omitempty"`
	ProductSetID       *string    `json:"product_set_id,omitempty"`
	Name               string     `json:"name"`
	StartAt            *time.Time `json:"start_at"`
	EndAt              *time.Time `json:"end_at"`
	Amount             int        `json:"amount"`
	TotalThreshold     *int       `json:"total_threshold,omitempty"`
	Type               string     `json:"type"`
	Target             string     `json:"target"`
	Created            time.Time  `json:"created"`
	Modified           time.Time  `json:"modified"`
}

// CreatePromoRule calls the API to create a new promo rule.
func (c *EcomClient) CreatePromoRule(ctx context.Context, p *PromoRuleRequest) (*PromoRule, error) {
	request, err := json.Marshal(&p)
	if err != nil {
		return nil, fmt.Errorf("%w: client: json marshal failed", err)
	}

	uri := c.endpoint + "/promo-rules"
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, uri, body)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("%w: client decode error", err)
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var promoRule PromoRule
	if err = json.NewDecoder(res.Body).Decode(&promoRule); err != nil {
		return nil, fmt.Errorf("%w: decode failed", err)
	}
	return &promoRule, nil
}

// GetPromoRule returns a single promo rule
func (c *EcomClient) GetPromoRule(ctx context.Context, promoRuleID string) (*PromoRule, error) {
	url := fmt.Sprintf("%s/promo-rules/%s", c.endpoint, promoRuleID)
	res, err := c.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}

		if e.Code == "promo-rules/promo-rule-not-found" {
			return nil, ErrPromoRuleNotFound
		}
		return nil, fmt.Errorf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
	}

	var p PromoRule
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &p, nil
}

// GetPromoRules returns a list of all promo rules.
func (c *EcomClient) GetPromoRules(ctx context.Context) ([]*PromoRule, error) {
	uri := c.endpoint + "/promo-rules"
	res, err := c.request(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer res.Body.Close()

	var container PromoRulesContainer
	if err := json.NewDecoder(res.Body).Decode(&container); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return container.Data, nil
}

// DeletePromoRule deletes a promo rule by id
func (c *EcomClient) DeletePromoRule(ctx context.Context, id string) error {
	uri := c.endpoint + "/promo-rules/" + id
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		return nil
	}
	if res.StatusCode == 400 {
		return ErrBadRequest
	}
	if res.StatusCode == 404 {
		return ErrPromoRuleNotFound
	}
	return nil
}
