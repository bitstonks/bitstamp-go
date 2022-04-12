package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// A helper function, custom URL merging logic adapted for the API.
func urlMerge(baseUrl url.URL, urlPath string, queryParams ...[2]string) string {
	baseUrl.Path = path.Join(baseUrl.Path, urlPath)

	// apparently, path.Join loses trailing slash in urlPath. we don't want that...
	if strings.HasSuffix(urlPath, "/") {
		baseUrl.Path += "/"
	}

	// add query params
	values := baseUrl.Query()
	for _, param := range queryParams {
		values.Set(param[0], param[1])
	}
	baseUrl.RawQuery = values.Encode()

	return baseUrl.String()
}

func validateCurrencyPair(currencyPair string) error {
	if _, exists := roundings[currencyPair]; exists {
		return nil
	} else {
		return fmt.Errorf("unknown currency pair: %s", currencyPair)
	}
}

// GetRequestError is a custom error type that makes for somewhat nicer logic with non-200 codes returned.
type GetRequestError struct {
	Code    int
	Content string
	Status  string
	Url     string
}

func (e *GetRequestError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Status, e.Url)
}

// HttpClient implements the HTTP (REST) API endpoints.
type HttpClient struct {
	*httpClientConfig
}

func NewHttpClient(options ...HttpOption) *HttpClient {
	config := defaultHttpClientConfig()
	for _, option := range options {
		option(config)
	}
	return &HttpClient{config}
}

func (c *HttpClient) credentials() url.Values {
	nonce := c.nonceGenerator()
	message := nonce + c.username + c.apiKey

	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(message))
	signature := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	data := make(url.Values)
	data.Set("key", c.apiKey)
	data.Set("signature", signature)
	data.Set("nonce", nonce)
	return data
}

func (c *HttpClient) getRequest(responseObject interface{}, urlPath string, queryParams ...[2]string) (err error) {
	url_ := urlMerge(c.domain, urlPath, queryParams...)

	resp, err := http.Get(url_)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return &GetRequestError{
			Code:    resp.StatusCode,
			Content: string(respBody),
			Status:  resp.Status,
			Url:     url_,
		}
	}

	err = json.Unmarshal(respBody, responseObject)
	return
}

func (c *HttpClient) authenticatedPostRequest(responseObject interface{}, urlPath string, queryParams ...[2]string) (err error) {
	authVersion := "v2"
	method := "POST"
	xAuth := "BITSTAMP " + c.apiKey
	apiSecret := []byte(c.apiSecret)
	contentType := "application/x-www-form-urlencoded"
	timestamp_ := c.timestampGenerator()
	nonce := c.nonceGenerator()
	url_ := urlMerge(c.domain, urlPath)

	var payloadString string
	if queryParams != nil {
		urlParams := url.Values{}
		for _, p := range queryParams {
			urlParams.Set(p[0], p[1]) // TODO: or is it .Add() here? any array arguments in the documentation?
		}
		payloadString = urlParams.Encode()
	}

	// message construction
	msg := xAuth + method + strings.TrimPrefix(strings.TrimPrefix(url_, "https://"), "http://")
	if queryParams == nil {
		msg = msg + nonce + timestamp_ + authVersion // TODO: apparently, contentType must be omitted here?
	} else {
		msg = msg + contentType + nonce + timestamp_ + authVersion + payloadString
	}
	sig := hmac.New(sha256.New, apiSecret)
	sig.Write([]byte(msg))
	signature := hex.EncodeToString(sig.Sum(nil))

	// do the request
	client := &http.Client{}
	var req *http.Request
	if queryParams == nil {
		req, err = http.NewRequest(method, url_, nil)
	} else {
		req, err = http.NewRequest(method, url_, bytes.NewBuffer([]byte(payloadString)))
	}
	if err != nil {
		return
	}
	req.Header.Add("X-Auth", xAuth)
	req.Header.Add("X-Auth-Signature", signature)
	req.Header.Add("X-Auth-Nonce", nonce)
	req.Header.Add("X-Auth-Timestamp", timestamp_)
	req.Header.Add("X-Auth-Version", authVersion)
	if queryParams != nil {
		req.Header.Add("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// handle response
	if resp.StatusCode != 200 {
		var errorMsg map[string]interface{}
		err = json.Unmarshal(respBody, &errorMsg)
		if err != nil {
			return
		}

		reasonVal, reasonPresent := errorMsg["reason"]
		codeVal, codePresent := errorMsg["code"]
		if reasonPresent && codePresent {
			err = fmt.Errorf("%s %s (%d)", codeVal, reasonVal, resp.StatusCode)
			return
		} else {
			err = fmt.Errorf("%s (%d)", string(respBody), resp.StatusCode)
			return
		}
	} else {
		// verify server signature
		checkMsg := nonce + timestamp_ + resp.Header.Get("Content-Type") + string(respBody)
		sig := hmac.New(sha256.New, apiSecret)
		sig.Write([]byte(checkMsg))
		serverSig := hex.EncodeToString(sig.Sum(nil))
		if serverSig != resp.Header.Get("X-Server-Auth-Signature") {
			err = fmt.Errorf("server signature mismatch: us (%s) them (%s)", serverSig, resp.Header.Get("X-Server-Auth-Signature"))
			return
		}

		err = json.Unmarshal(respBody, responseObject)
		if err != nil {
			return
		}
	}

	return
}
