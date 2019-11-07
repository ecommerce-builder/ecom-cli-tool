package eclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ImageRequest JSON request body.
type ImageRequest struct {
	ProductID string `json:"product_id"`
	Path      string `json:"path"`
}

// ImageResponse JSON image response body.
type ImageResponse struct {
	Object      string    `json:"object"`
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	ProductPath string    `json:"product_path"`
	ProducutSKU string    `json:"product_sku"`
	Path        string    `json:"path"`
	GSURL       string    `json:"gsurl"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Size        int       `json:"size"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

// CreateImage calls the API service to create an image for a given product.
func (c *EcomClient) CreateImage(image ImageRequest) (*ImageResponse, error) {
	request, err := json.Marshal(&image)
	if err != nil {
		return nil, fmt.Errorf("client: json marshal failed: %w", err)
	}

	params := url.Values{}
	params.Add("product_id", image.ProductID)
	uri := c.endpoint + "/images?" + params.Encode()
	body := strings.NewReader(string(request))
	res, err := c.request(http.MethodPost, uri, body)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("client decode error: %w", err)
		}
		return nil, fmt.Errorf("Status: %d, Code: %s, Message: %s", e.Status, e.Code, e.Message)
	}

	var response ImageResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("response decode failed: %w", err)
	}
	return &response, nil
}

// DeleteProductImages calls the API Service to delete all
// images for a given product id.
func (c *EcomClient) DeleteProductImages(productID string) error {
	params := url.Values{}
	params.Add("product_id", productID)

	uri := c.endpoint + "/images?" + params.Encode()
	res, err := c.request(http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("client decode error: %w", err)
		}
		return fmt.Errorf("Status: %d, Code: %s, Message: %s", e.Status, e.Code, e.Message)
	}
	return nil
}
