package gplay

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type requestOptions struct {
	Method  string
	URL     string
	Body    []byte
	Headers http.Header
}

type RequestError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *RequestError) Error() string {
	if e == nil {
		return ""
	}
	if e.StatusCode != 0 {
		return fmt.Sprintf("%s (status=%d)", e.Message, e.StatusCode)
	}
	return e.Message
}

func (e *RequestError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (c *Client) doRequest(ctx context.Context, opts requestOptions, throttlePerSecond int) ([]byte, int, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	method := opts.Method
	if method == "" {
		method = http.MethodGet
	}

	u, err := url.Parse(c.withBaseURL(opts.URL))
	if err != nil {
		return nil, 0, err
	}

	retryCount := 0
	retryWait := 0 * time.Millisecond
	if c != nil {
		retryCount = c.retryCount
		retryWait = c.retryWait
	}

	attempts := retryCount + 1
	if attempts < 1 {
		attempts = 1
	}

	for attempt := 0; attempt < attempts; attempt++ {
		if err := c.throttle.wait(ctx, throttlePerSecond); err != nil {
			return nil, 0, err
		}

		var body io.Reader
		if len(opts.Body) > 0 {
			body = bytes.NewReader(opts.Body)
		}

		req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
		if err != nil {
			return nil, 0, err
		}
		for k, vv := range opts.Headers {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; google-play-scraper-go)")
		}
		if req.Header.Get("Accept-Language") == "" {
			req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt < attempts-1 {
				if err := sleepCtx(ctx, backoff(retryWait, attempt)); err != nil {
					return nil, 0, err
				}
				continue
			}
			return nil, 0, &RequestError{Message: "Error requesting Google Play: " + err.Error(), Err: err}
		}

		b, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			if attempt < attempts-1 {
				if err := sleepCtx(ctx, backoff(retryWait, attempt)); err != nil {
					return nil, 0, err
				}
				continue
			}
			return nil, resp.StatusCode, readErr
		}

		if resp.StatusCode >= 400 {
			retryable := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable
			if retryable && attempt < attempts-1 {
				if err := sleepCtx(ctx, backoff(retryWait, attempt)); err != nil {
					return nil, resp.StatusCode, err
				}
				continue
			}
			msg := "Error requesting Google Play: " + resp.Status
			if resp.StatusCode == http.StatusNotFound {
				msg = "App not found (404)"
			}
			return nil, resp.StatusCode, &RequestError{StatusCode: resp.StatusCode, Message: msg, Err: errors.New(string(b))}
		}

		return b, resp.StatusCode, nil
	}

	return nil, 0, errors.New("request failed")
}

func backoff(base time.Duration, attempt int) time.Duration {
	if base <= 0 {
		base = 400 * time.Millisecond
	}
	mult := 1
	for i := 0; i < attempt; i++ {
		mult *= 2
		if mult > 16 {
			mult = 16
			break
		}
	}
	return time.Duration(mult) * base
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
