package eclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// EcomClient structure.
type EcomClient struct {
	firebaseAPIKey string
	endpoint       string
	client         *http.Client
	jwt            string
}

type sysInfoPg struct {
	Host        string `json:"ECOM_PG_HOST"`
	Port        string `json:"ECOM_PG_PORT"`
	Database    string `json:"ECOM_PG_DATABASE"`
	User        string `json:"ECOM_PG_USER"`
	SslMode     string `json:"ECOM_PG_SSLMODE"`
	SslCert     string `json:"ECOM_PG_SSLCERT"`
	SslKey      string `json:"ECOM_PG_SSLKEY"`
	SslRootCert string `json:"ECOM_PG_SSLROOTCERT"`
}

type sysInfoGoog struct {
	ProjectID string `json:"ECOM_GOOGLE_PROJECT_ID"`
}

type sysInfoApp struct {
	HTTPPort  string `json:"PORT"`
	RootEmail string `json:"ECOM_APP_ROOT_EMAIL"`
}

type sysInfoEnv struct {
	Pg   sysInfoPg   `json:"pg"`
	Goog sysInfoGoog `json:"google"`
	App  sysInfoApp  `json:"app"`
}

// SysInfo provides a record of system information.
type SysInfo struct {
	APIVersion string     `json:"api_version"`
	Env        sysInfoEnv `json:"env"`
}

type devKeyRequest struct {
	Key string `json:"key"`
}

type customToken struct {
	CustomToken string `json:"custom_token"`
}

// NewEcomClient creates an EcomClient struct for interacting with the API Service
func NewEcomClient(firebaseAPIKey, endpoint string, timeout time.Duration) *EcomClient {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 10,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	return &EcomClient{
		endpoint:       endpoint,
		client:         client,
		firebaseAPIKey: firebaseAPIKey,
	}
}

// SetJWT sets the current Firebase JWT for future calls to the e-commerce API.
func (c *EcomClient) SetJWT(jwt string) {
	c.jwt = jwt
}

// URLToHostName converts a standard URL string to hostname replacing the dot character with underscores.
func URLToHostName(u string) (string, error) {
	url, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse url %q", u)
	}
	return strings.ReplaceAll(url.Hostname(), ".", "_"), nil
}

// TokenAndRefreshToken contains the latest token (1hr expiration) and refresh token (last indefinitely).
type TokenAndRefreshToken struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

// See https://firebase.google.com/docs/reference/rest/auth/#section-verify-custom-token
// token	      string   A Firebase Auth custom token from which to create an ID and refresh token pair.
// returnSecureToken  boolean  Whether or not to return an ID and refresh token. Should always be true.
type verifyCustomTokenRequest struct {
	Token             string `json:"token"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type verifyCustomTokenResponse struct {
	Kind         string `json:"kind"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

// curl -H 'Content-Type: application/x-www-form-urlencoded' -X POST --data 'grant_type=refresh_token&refresh_token=AEu4IL1BsyHyQ7lfBaUXrukvZfOJ5KEOjTYpMueRimrPmQ00GioTbIsAPsuWAG6JEp5o2SBVBpNCySu3OsxBFstDbPaQnrGYKUtRMw9ENqTt7Qmq9Sdy7LzNkxu7cizlxiq2bDVuj80DAmh1oUP_rjehBUMk1HUO4UtN737Ggk1IGHFf4-rTxCZtF5nUoqO8W34S53Ik32RdK3QvbRdRlwav_xwiXyM0UA' https://securetoken.googleapis.com/v1/token?key=AIzaSyBGU4AnEHCOXKGkOtwXWyxOBaU3VSTg6wY

// grant_type	  string  The refresh token's grant type, always "refresh_token".
// refresh_token  string  A Firebase Auth refresh token.
type exchangeRefreshTokenRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type exchangeRefreshTokenResponse struct {
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
	ProjectID    string `json:"project_id"`
}

type ertBadRequestResponse struct {
	Error struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// ExchangeRefreshTokenForIDToken calls Google's REST API.
// Response Payload
// Property Name	Type	Description
// expires_in	string	The number of seconds in which the ID token expires.
// token_type	string	The type of the refresh token, always "Bearer".
// refresh_token	string	The Firebase Auth refresh token provided in the request or a new refresh token.
// id_token	string	A Firebase Auth ID token.
// user_id	string	The uid corresponding to the provided ID token.
// project_id	string	Your Firebase project ID.
func (c *EcomClient) ExchangeRefreshTokenForIDToken(refreshToken string) (*TokenAndRefreshToken, error) {
	v := url.Values{}
	v.Set("key", c.firebaseAPIKey)
	uri := url.URL{
		Scheme:     "https",
		Host:       "securetoken.googleapis.com",
		Path:       "v1/token",
		ForceQuery: false,
		RawQuery:   v.Encode(),
	}

	reqBody := exchangeRefreshTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	}
	payload := url.Values{}
	payload.Set("grant_type", reqBody.GrantType)
	payload.Set("refresh_token", reqBody.RefreshToken)
	req, err := http.NewRequest("POST", uri.String(), strings.NewReader(payload.Encode()))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "create new POST request failed")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var badReqRes ertBadRequestResponse
		body, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &badReqRes)
		if err != nil {
			return nil, errors.Wrapf(err, "%d %s\n", badReqRes.Error.Code, badReqRes.Error.Message)
		}

	}

	response := exchangeRefreshTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "json ecode error")
	}
	return &TokenAndRefreshToken{
		IDToken:      response.IDToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

// ExchangeCustomTokenForIDAndRefreshToken calls the Firebase REST API to exchange a customer token for Firebase token and refresh token.
func (c *EcomClient) ExchangeCustomTokenForIDAndRefreshToken(customToken string) (*TokenAndRefreshToken, error) {
	// build the URL including Query params
	v := url.Values{}
	v.Set("key", c.firebaseAPIKey)
	uri := url.URL{
		Scheme:     "https",
		Host:       "www.googleapis.com",
		Path:       "identitytoolkit/v3/relyingparty/verifyCustomToken",
		ForceQuery: false,
		RawQuery:   v.Encode(),
	}

	// build and execute the request
	reqBody := verifyCustomTokenRequest{
		Token:             customToken,
		ReturnSecureToken: true,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(reqBody)
	req, err := http.NewRequest("POST", uri.String(), buf)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "creating new POST request failed")
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		var badReqRes struct {
			Error struct {
				Code    int64  `json:"code"`
				Message string `json:"message"`
				Errors  []struct {
					Message string `json:"message"`
					Domain  string `json:"domain"`
					Reason  string `json:"reason"`
				} `json:"errors"`
				Status string `json:"status"`
			} `json:"error"`
		}
		err = json.NewDecoder(res.Body).Decode(&badReqRes)
		if err != nil {
			return nil, errors.Wrap(err, "decode bad request response failed")
		}
		return nil, fmt.Errorf("%d %s", badReqRes.Error.Code, badReqRes.Error.Message)
	} else if res.StatusCode > 400 {
		return nil, fmt.Errorf("%s", res.Status)
	}

	tokenResponse := verifyCustomTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, errors.Wrap(err, "json decode failed")
	}
	return &TokenAndRefreshToken{
		IDToken:      tokenResponse.IDToken,
		RefreshToken: tokenResponse.RefreshToken,
	}, nil
}

// https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=[API_KEY]

// SignInWithDevKey exchanges a Developer Key for a Custom Token.
func (c *EcomClient) SignInWithDevKey(key string) (token string, err error) {
	uri := c.endpoint + "/signin-with-devkey"
	payload := devKeyRequest{
		Key: key,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)
	req, err := http.NewRequest("POST", uri, buf)
	if err != nil {
		return "", fmt.Errorf("error creating new POST request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing HTTP POST to %v : %v", uri, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", errors.Wrapf(err, "%s", res.Status)
	}

	ct := customToken{}
	err = json.NewDecoder(res.Body).Decode(&ct)
	if err != nil {
		return "", errors.Wrap(err, "custom token json decode error")
	}
	return ct.CustomToken, nil
}

// SysInfo retrieves the System Info from the API endpoint.
func (c *EcomClient) SysInfo() (*SysInfo, error) {
	uri := c.endpoint + "/sysinfo"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating new GET request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.jwt)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing HTTP GET to %v : %v", uri, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, errors.Wrapf(err, "%s", res.Status)
	}

	var sysInfo SysInfo
	if err := json.NewDecoder(res.Body).Decode(&sysInfo); err != nil {
		return nil, errors.Wrapf(err, "failed to decode url %s", uri)
	}
	return &sysInfo, nil
}