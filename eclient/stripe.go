package eclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type stripeCheckoutResponseBody struct {
	Object            string `json:"object"`
	CheckoutSessionID string `json:"checkout_session_id"`
}

// StripeCheckout calls the API service to generate a stripe checkout session
// for a given order.
func (c *EcomClient) StripeCheckout(ctx context.Context, orderID string) (string, error) {
	url := fmt.Sprintf("%s/orders/%s/stripecheckout", c.endpoint, orderID)
	res, err := c.request(http.MethodPost, url, nil)
	if err != nil {
		return "", errors.Wrapf(err, "c.request(http.MethodPost, url=%q, nil)", url)
	}

	if res.StatusCode >= 400 {
		var e badRequestResponse
		dec := json.NewDecoder(res.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&e); err != nil {
			return "", errors.Wrap(err, "decode")
		}
		return "", fmt.Errorf("status: %d, code: %s, message: %s",
			e.Status, e.Code, e.Message)
	}

	var v stripeCheckoutResponseBody
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return "", errors.Wrap(err, "decode")
	}
	return v.CheckoutSessionID, nil
}
