package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// A helper function, custom URL merging logic adapted for the API.
func urlMerge(baseUrl url.URL, urlPath string, queryParams *url.Values) string {
	baseUrl.Path = path.Join(baseUrl.Path, urlPath)

	// apparently, path.Join loses trailing slash in urlPath. we don't want that...
	if strings.HasSuffix(urlPath, "/") {
		baseUrl.Path += "/"
	}

	// add query params
	values := baseUrl.Query()
	if queryParams != nil {
		for a, param := range *queryParams {
			for _, value := range param {
				values.Add(a, value)
			}
		}
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

func (c *HttpClient) getRequest(responseObject interface{}, urlPath string, queryParams *url.Values) (err error) {
	url_ := urlMerge(c.domain, urlPath, queryParams)

	resp, err := http.Get(url_)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
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

func (c *HttpClient) authenticatedFormRequest(responseObject interface{}, method string, urlPath string, queryParams *url.Values, formQueryParams map[string]string) (err error) {
	contentType := "application/x-www-form-urlencoded"
	var payloadString string
	if formQueryParams != nil {
		urlParams := url.Values{}
		for paramName, paramVal := range formQueryParams {
			urlParams.Add(paramName, paramVal)
		}
		payloadString = urlParams.Encode()
	}

	err = c.doSignedRequest(responseObject, method, urlPath, queryParams, contentType, payloadString)
	return
}

func (c *HttpClient) authenticatedJsonRequest(responseObject interface{}, method string, urlPath string, urlParams *url.Values, requestObject interface{}) (err error) {
	contentType := "application/json"
	var payloadString string
	var payloadBytes []byte
	if requestObject != nil {
		payloadBytes, err = json.Marshal(requestObject)
		if payloadBytes != nil {
			payloadString = string(payloadBytes)
		}
	}

	err = c.doSignedRequest(responseObject, method, urlPath, urlParams, contentType, payloadString)
	return
}

type PaginationWrapper struct {
	Data interface{} `json:"data"`
}

func (c *HttpClient) doSignedRequest(responseObject interface{}, method string, urlPath string, urlParams *url.Values, contentType string, payloadString string) (err error) {
	url_ := urlMerge(c.domain, urlPath, urlParams)
	authVersion := "v2"
	xAuth := "BITSTAMP " + c.apiKey
	apiSecret := []byte(c.apiSecret)
	timestamp_ := c.timestampGenerator()
	nonce := c.nonceGenerator()
	// message construction
	msg := xAuth + method + strings.TrimPrefix(strings.TrimPrefix(url_, "https://"), "http://")
	if payloadString == "" {
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
	if payloadString == "" {
		req, err = http.NewRequest(method, url_, nil)
	} else {
		req, err = http.NewRequest(method, url_, bytes.NewBuffer([]byte(payloadString)))
	}
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth", xAuth)
	req.Header.Add("X-Auth-Signature", signature)
	req.Header.Add("X-Auth-Nonce", nonce)
	req.Header.Add("X-Auth-Timestamp", timestamp_)
	req.Header.Add("X-Auth-Version", authVersion)
	if payloadString != "" {
		req.Header.Add("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// handle response
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 204 {
		if resp.StatusCode == 503 {
			err = errors.New("service unavailable")
			return
		}
		var errorMsg map[string]interface{}
		err = json.Unmarshal(respBody, &errorMsg)
		if err != nil {
			return err
		}

		reasonVal, reasonPresent := errorMsg["reason"]
		codeVal, codePresent := errorMsg["code"]
		if reasonPresent && codePresent {
			err = fmt.Errorf("%s %s (%d)", codeVal, reasonVal, resp.StatusCode)
			return err
		} else {
			err = fmt.Errorf("%s (%d)", string(respBody), resp.StatusCode)
			return err
		}
	} else {
		// verify server signature
		checkMsg := nonce + timestamp_ + resp.Header.Get("Content-Type") + string(respBody)
		sig := hmac.New(sha256.New, apiSecret)
		sig.Write([]byte(checkMsg))
		serverSig := hex.EncodeToString(sig.Sum(nil))
		if serverSig != resp.Header.Get("X-Server-Auth-Signature") {
			err = fmt.Errorf("server signature mismatch: us (%s) them (%s)", serverSig, resp.Header.Get("X-Server-Auth-Signature"))
			return err
		}
		if len(respBody) > 0 {
			err = json.Unmarshal(respBody, responseObject)
			if err != nil {
				var wrapped PaginationWrapper
				wrapped = PaginationWrapper{Data: responseObject}
				err = json.Unmarshal(respBody, &wrapped)
				responseObject = wrapped.Data
				if err != nil {
					return err

				}
			}
		}
	}

	return nil
}
