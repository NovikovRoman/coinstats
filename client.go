package coinstats

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const host = "https://openapiv1.coinstats.app"

type Client struct {
	apiKey     string
	shareToken string
	transport  *http.Transport
	httpClient *http.Client
}

type Opt func(*Client)

func WithTransport(transport *http.Transport) Opt {
	return func(c *Client) {
		c.transport = transport
	}
}

func WithShareToken(token string) Opt {
	return func(c *Client) {
		c.shareToken = token
	}
}

func New(apiKey string, opts ...Opt) *Client {
	c := &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		transport: &http.Transport{
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}
	c.httpClient.Transport = c.transport
	return c
}

func (c *Client) addApiKey(req *http.Request) {
	req.Header.Set("X-API-KEY", c.apiKey)
}

func (c *Client) hasShareToken() bool {
	return c.shareToken != ""
}

func (c *Client) addShareToken(req *http.Request) {
	if c.hasShareToken() {
		req.Header.Set("sharetoken", c.shareToken)
	}
}

func do[T any](client *Client, req *http.Request) (T, error) {
	client.addApiKey(req)

	var data T
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return data, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}()

	b, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case 400:
		err = ErrBadRequest
	case 401:
		err = ErrUnauthorized
	case 403:
		err = ErrForbidden
	case 404:
		err = ErrNotFound
	case 409:
		err = ErrTransactionsNotSynced
	case 429:
		err = ErrRateLimit
	case 503:
		err = ErrServiceUnavailable
	case 200:
		err = nil
	default:
		err = ErrBadStatusCode
		err = errors.Join(err, fmt.Errorf("%d %s", resp.StatusCode, string(b)))
		return data, err
	}

	if err != nil {
		var e Error
		jsonErr := json.Unmarshal(b, &e)
		if jsonErr != nil {
			err = errors.Join(err, fmt.Errorf("%w %s", jsonErr, string(b)))
		} else {
			err = errors.Join(err, &e)
		}
		return data, err
	}

	err = json.Unmarshal(b, &data)
	return data, err
}
