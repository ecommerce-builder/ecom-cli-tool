package eclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ecommerce-builder/ecom-cli-tool/configmgr"
	"github.com/pkg/errors"
)

// Version string
var Version string

var userAgent string

// EcomClient structure.
type EcomClient struct {
	endpoint string
	scheme   string
	hostname string
	port     string
	client   *http.Client
	jwt      string
}

type sysInfoPg struct {
	Host          string `json:"ECOM_PG_HOST"`
	Port          string `json:"ECOM_PG_PORT"`
	Database      string `json:"ECOM_PG_DATABASE"`
	User          string `json:"ECOM_PG_USER"`
	SslMode       string `json:"ECOM_PG_SSLMODE"`
	SslCert       string `json:"ECOM_PG_SSLCERT"`
	SslKey        string `json:"ECOM_PG_SSLKEY"`
	SslRootCert   string `json:"ECOM_PG_SSLROOTCERT"`
	SchemaVersion string `json:"schema_version"`
}

type sysInfoGoog struct {
	ProjectID string `json:"ECOM_GAE_PROJECT_ID"`
}

type sysInfoFirebase struct {
	APIKEY            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	DatabaseURL       string `json:"databaseURL"`
	ProjectID         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderID string `json:"messagingSenderId"`
	AppID             string `json:"appId"`
}

type sysInfoStripe struct {
	StripeSuccessURL string `json:"ECOM_STRIPE_SUCCESS_URL"`
	StripeCancelURL  string `json:"ECOM_STRIPE_CANCEL_URL"`
}

type sysInfoApp struct {
	AppPort                     string `json:"PORT"`
	AppRootEmail                string `json:"ECOM_APP_ROOT_EMAIL"`
	AppEnableStackDriverLogging bool   `json:"ECOM_APP_ENABLE_STACKDRIVER_LOGGING"`
	AppEndpoint                 string `json:"ECOM_APP_ENDPOINT"`
}

type sysInfoEnv struct {
	Pg       sysInfoPg       `json:"pg"`
	Goog     sysInfoGoog     `json:"google"`
	Firebase sysInfoFirebase `json:"firebase"`
	Stripe   sysInfoStripe   `json:"stripe"`
	App      sysInfoApp      `json:"app"`
}

// SysInfo provides a record of system information.
type SysInfo struct {
	APIVersion string     `json:"api_version"`
	Env        sysInfoEnv `json:"env"`
}

type devKeyRequest struct {
	Key string `json:"key"`
}

type tokenAndCustomerResponse struct {
	CustomToken string       `json:"custom_token"`
	User        UserResponse `json:"user"`
}

var timeout = time.Duration(10 * time.Second)

// New creates an EcomClient struct for interacting with the API Service
func New(endpoint string) *EcomClient {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 10,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	return &EcomClient{
		endpoint: endpoint,
		scheme: url.Scheme,
		hostname: url.Host,
		port: url.Port(),
		client:   client,
	}
}

// SetJWT sets the current Firebase JWT for future calls to the e-commerce API.
func (c *EcomClient) SetJWT(jwt string) {
	c.jwt = jwt
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

type badRequestResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SetToken accepts an EcomConfigEntry and derives the token and refresh
// token file, before reading it, inspecting it and if necessary generating
// a refresh token, before writing back the file. The token is then stored
// in the EcomClient struct.
func (c *EcomClient) SetToken(cfg *configmgr.EcomConfigEntry) error {
	file, err := configmgr.TokenFilename(cfg)
	if err != nil {
		return fmt.Errorf("token file %q not found: %w", file, err)
	}
	tar, err := configmgr.ReadTokenAndRefreshToken(file)
	if err != nil {
		return fmt.Errorf("token and refresh token cannot be read from %q: %w", file, err)
	}
	var p jwt.Parser
	t, _, err := p.ParseUnverified(tar.IDToken, &jwt.StandardClaims{})
	claims := t.Claims.(*jwt.StandardClaims)
	utcNow := time.Now().Unix()

	// If the token has expired, use the refresh token to get another
	if claims.ExpiresAt-utcNow <= 0 {
		f, err := c.GetConfig()
		if err != nil {
			return err
		}
		tar, err = c.ExchangeRefreshTokenForIDToken(f.APIKEY, tar.RefreshToken)
		if err != nil {
			return fmt.Errorf("exchange refresh token for id token failed: %w", err)
		}
		hostname, err := configmgr.URLToHostName(cfg.Endpoint)
		filename := fmt.Sprintf("%s-%s", hostname, cfg.DevKey[:6])
		err = configmgr.WriteTokenAndRefreshToken(filename, tar)
		if err != nil {
			return fmt.Errorf("write token and refresh token failed: %w", err)
		}
	}
	c.jwt = tar.IDToken
	return nil
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
func (c *EcomClient) ExchangeRefreshTokenForIDToken(firebaseAPIKey, refreshToken string) (*configmgr.TokenAndRefreshToken, error) {
	v := url.Values{}
	v.Set("key", firebaseAPIKey)
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
		return nil, fmt.Errorf("create new POST request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var e badRequestResponse
		body, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &e)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", e.Code, e.Message, err)
		}
	}
	response := exchangeRefreshTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("json ecode error: %w", err)
	}
	return &configmgr.TokenAndRefreshToken{
		IDToken:      response.IDToken,
		RefreshToken: response.RefreshToken,
	}, nil
}

// ExchangeCustomTokenForIDAndRefreshToken calls the Firebase REST API to exchange a customer token for Firebase token and refresh token.
func (c *EcomClient) ExchangeCustomTokenForIDAndRefreshToken(firebaseAPIKey, token string) (*configmgr.TokenAndRefreshToken, error) {
	// build the URL including Query params
	v := url.Values{}
	v.Set("key", firebaseAPIKey)
	uri := url.URL{
		Scheme:     "https",
		Host:       "www.googleapis.com",
		Path:       "identitytoolkit/v3/relyingparty/verifyCustomToken",
		ForceQuery: false,
		RawQuery:   v.Encode(),
	}

	// build and execute the request
	reqBody := verifyCustomTokenRequest{
		Token:             token,
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
		return nil, fmt.Errorf("creating new POST request failed: %w", err)
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
			return nil, fmt.Errorf("decode failed: %w", err)
		}
		return nil, fmt.Errorf("%d %s", badReqRes.Error.Code, badReqRes.Error.Message)
	} else if res.StatusCode > 400 {
		return nil, fmt.Errorf("%s", res.Status)
	}

	tokenResponse := verifyCustomTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}
	return &configmgr.TokenAndRefreshToken{
		IDToken:      tokenResponse.IDToken,
		RefreshToken: tokenResponse.RefreshToken,
	}, nil
}

// SignInWithDevKey exchanges a Developer Key for a Customer token.
// https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=[API_KEY]
func (c *EcomClient) SignInWithDevKey(key string) (token string, user *UserResponse, err error) {
	uri := c.endpoint + "/signin-with-devkey"
	payload := devKeyRequest{
		Key: key,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)

	req, err := http.NewRequest("POST", uri, buf)
	if err != nil {
		return "", nil, fmt.Errorf("error creating new POST request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("error executing HTTP POST to %v : %w", uri, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", nil, fmt.Errorf("%s", res.Status)
	}

	ct := tokenAndCustomerResponse{}
	err = json.NewDecoder(res.Body).Decode(&ct)
	if err != nil {
		return "", nil, fmt.Errorf("custom token json decode error: %w", err)
	}
	return ct.CustomToken, &ct.User, nil
}

// SysInfo retrieves the System Info from the API endpoint.
func (c *EcomClient) SysInfo() (*SysInfo, error) {
	uri := c.endpoint + "/sysinfo"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("http new request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.jwt)
	req.Header.Set("User-Agent", fmt.Sprintf("ecom/%s", Version))
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET to %v failed: %w", uri, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s", res.Status)
	}

	var sysInfo SysInfo
	if err := json.NewDecoder(res.Body).Decode(&sysInfo); err != nil {
		return nil, fmt.Errorf("failed to decode url %s: %w", uri, err)
	}
	return &sysInfo, nil
}

// ConfigContainerResponse container
type ConfigContainerResponse struct {
	Object         string                  `json:"object"`
	FirebaseConfig *FirebaseConfigResponse `json:"firebaseConfig"`
}

// FirebaseConfigResponse firebase config response
type FirebaseConfigResponse struct {
	APIKEY            string `json:"apiKey"`
	AuthDomain        string `json:"authDomain"`
	DatabaseURL       string `json:"databaseURL"`
	ProjectID         string `json:"projectId"`
	StorageBucket     string `json:"storageBucket"`
	MessagingSenderID string `json:"messagingSenderId"`
	AppID             string `json:"appId"`
}

// GetConfig gets the Firebase Config from the server.
// HTTP GET /config is a public resource and requires no
// authorization or token.
func (c *EcomClient) GetConfig() (*FirebaseConfigResponse, error) {
	uri := c.endpoint + "/config"
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("http new request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do to %v failed: %w", uri, err)
	}
	defer res.Body.Close()
	var f *ConfigContainerResponse
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		return nil, errors.Wrapf(err, "json decode url %s", uri)
	}
	return f.FirebaseConfig, nil
}
