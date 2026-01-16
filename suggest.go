package gplay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var suggestURLTemplate = "/_/PlayStoreUi/data/batchexecute?rpcids=IJ4APc&f.sid=-697906427155521722&bl=boq_playuiserver_20190903.08_p0&hl=%s&gl=%s&authuser&soc-app=121&soc-platform=1&soc-device=1&_reqid=1065213"

func (c *Client) Suggest(ctx context.Context, opts SuggestOptions) ([]string, error) {
	if opts.Term == "" {
		return nil, errors.New("term missing")
	}
	lang := opts.Lang
	if lang == "" {
		lang = "en"
	}
	country := opts.Country
	if country == "" {
		country = "us"
	}
	cacheOpts := opts
	cacheOpts.Lang = lang
	cacheOpts.Country = country
	if c != nil && c.cache != nil {
		var cached []string
		hit, err := c.cacheGet("suggest", cacheOpts, &cached)
		if err != nil {
			return nil, err
		}
		if hit {
			return cached, nil
		}
	}

	u := fmt.Sprintf(suggestURLTemplate, queryEscape(lang), queryEscape(country))
	term := queryEscape(opts.Term)
	body := fmt.Sprintf("f.req=%%5B%%5B%%5B%%22IJ4APc%%22%%2C%%22%%5B%%5Bnull%%2C%%5B%%5C%%22%s%%5C%%22%%5D%%2C%%5B10%%5D%%2C%%5B2%%5D%%2C4%%5D%%5D%%22%%5D%%5D%%5D", term)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	for k, vv := range opts.Headers {
		for _, v := range vv {
			headers.Add(k, v)
		}
	}

	respBody, _, err := c.do(ctx, requestOptions{Method: http.MethodPost, URL: u, Body: []byte(body), Headers: headers}, opts.Throttle)
	if err != nil {
		return nil, err
	}
	outer, err := parseBatchedExecuteResponse(respBody)
	if err != nil {
		return nil, err
	}
	inner, err := parseBatchedInnerJSON(outer)
	if err != nil {
		return nil, err
	}
	if inner == nil {
		return []string{}, nil
	}

	b, err := json.Marshal(inner)
	if err != nil {
		return nil, err
	}
	var data []any
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	root, ok := pathGet(data, []any{0, 0}).([]any)
	if !ok {
		return []string{}, nil
	}
	out := make([]string, 0, len(root))
	for _, it := range root {
		if s, ok := pathGet(it, []any{0}).(string); ok {
			out = append(out, s)
		}
	}
	c.cacheSet("suggest", cacheOpts, out)
	return out, nil
}
