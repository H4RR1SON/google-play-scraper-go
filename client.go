package gplay

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	throttle   *throttleState
	cache      *memoCache
	retryCount int
	retryWait  time.Duration
}

type ClientOptions struct {
	HTTPClient *http.Client
	BaseURL    string
	Timeout    time.Duration
	RetryCount int
	RetryWait  time.Duration
	ProxyURL   string
}

func NewClient(opts ClientOptions) (*Client, error) {
	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = BaseURL
	}

	hc := opts.HTTPClient
	if hc == nil {
		jar, _ := cookiejar.New(nil)
		hc = &http.Client{Jar: jar}
	}
	if hc.Jar == nil {
		jar, _ := cookiejar.New(nil)
		hc.Jar = jar
	}
	if opts.Timeout > 0 {
		hc.Timeout = opts.Timeout
	}
	if opts.ProxyURL != "" {
		if hc.Transport == nil {
			proxyURL, err := url.Parse(opts.ProxyURL)
			if err != nil {
				return nil, err
			}
			hc.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}

	retryWait := opts.RetryWait
	if retryWait == 0 {
		retryWait = 400 * time.Millisecond
	}

	return &Client{
		httpClient: hc,
		baseURL:    baseURL,
		throttle:   newThrottleState(),
		retryCount: opts.RetryCount,
		retryWait:  retryWait,
	}, nil
}

func MustNewClient(opts ClientOptions) *Client {
	c, err := NewClient(opts)
	if err != nil {
		panic(err)
	}
	return c
}

var DefaultClient = MustNewClient(ClientOptions{Timeout: 15 * time.Second})

func (c *Client) withBaseURL(url string) string {
	if c == nil {
		return DefaultClient.withBaseURL(url)
	}
	if len(url) >= 4 && url[:4] == "http" {
		return url
	}
	return c.baseURL + url
}

func (c *Client) do(ctx context.Context, reqOpts requestOptions, throttlePerSecond int) ([]byte, int, error) {
	if c == nil {
		return DefaultClient.do(ctx, reqOpts, throttlePerSecond)
	}
	return c.doRequest(ctx, reqOpts, throttlePerSecond)
}
