package gplay

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

func (c *Client) Developer(ctx context.Context, opts DeveloperOptions) ([]App, error) {
	if opts.DevID == "" {
		return nil, errors.New("devId missing")
	}
	lang := opts.Lang
	if lang == "" {
		lang = "en"
	}
	country := opts.Country
	if country == "" {
		country = "us"
	}
	num := opts.Num
	if num == 0 {
		num = 60
	}
	cacheOpts := opts
	cacheOpts.Lang = lang
	cacheOpts.Country = country
	cacheOpts.Num = num
	if c != nil && c.cache != nil {
		var cached []App
		hit, err := c.cacheGet("developer", cacheOpts, &cached)
		if err != nil {
			return nil, err
		}
		if hit {
			return cached, nil
		}
	}

	path := "/store/apps/developer"
	if _, err := strconv.ParseInt(opts.DevID, 10, 64); err == nil {
		path = "/store/apps/dev"
	}

	qs := url.Values{}
	qs.Set("id", opts.DevID)
	qs.Set("hl", lang)
	qs.Set("gl", country)
	pageURL := path + "?" + encodeValues(qs)

	body, _, err := c.do(ctx, requestOptions{URL: pageURL, Headers: opts.Headers}, opts.Throttle)
	if err != nil {
		return nil, err
	}

	parsed := parseScriptData(body)

	appsRoot := []any{"ds:3", 0, 1, 0, 22, 0}
	tokenPath := []any{"ds:3", 0, 1, 0, 22, 1, 3, 1}
	m := map[string]fieldSpec{
		"title":     {Path: []any{0, 3}},
		"appId":     {Path: []any{0, 0, 0}},
		"url":       {Path: []any{0, 10, 4, 2}, Fn: func(input any, _ parsedData) any { p, _ := asString(input); return resolveURL(BaseURL, p) }},
		"icon":      {Path: []any{0, 1, 3, 2}},
		"developer": {Path: []any{0, 14}},
		"currency":  {Path: []any{0, 8, 1, 0, 1}},
		"price":     {Path: []any{0, 8, 1, 0, 0}, Fn: func(input any, _ parsedData) any { v, _ := asFloat(input); return v / 1000000 }},
		"free":      {Path: []any{0, 8, 1, 0, 0}, Fn: func(input any, _ parsedData) any { v, _ := asFloat(input); return v == 0 }},
		"summary":   {Path: []any{0, 13, 1}},
		"scoreText": {Path: []any{0, 4, 0}},
		"score":     {Path: []any{0, 4, 1}},
	}

	if _, err := strconv.ParseInt(opts.DevID, 10, 64); err == nil {
		appsRoot = []any{"ds:3", 0, 1, 0, 21, 0}
		tokenPath = []any{"ds:3", 0, 1, 0, 21, 1, 3, 1}
		m = map[string]fieldSpec{
			"title":     {Path: []any{3}},
			"appId":     {Path: []any{0, 0}},
			"url":       {Path: []any{10, 4, 2}, Fn: func(input any, _ parsedData) any { p, _ := asString(input); return resolveURL(BaseURL, p) }},
			"icon":      {Path: []any{1, 3, 2}},
			"developer": {Path: []any{14}},
			"currency":  {Path: []any{8, 1, 0, 1}},
			"price":     {Path: []any{8, 1, 0, 0}, Fn: func(input any, _ parsedData) any { v, _ := asFloat(input); return v / 1000000 }},
			"free":      {Path: []any{8, 1, 0, 0}, Fn: func(input any, _ parsedData) any { v, _ := asFloat(input); return v == 0 }},
			"summary":   {Path: []any{13, 1}},
			"scoreText": {Path: []any{4, 0}},
			"score":     {Path: []any{4, 1}},
		}
	}

	appsAny := pathGet(parsed, appsRoot)
	appsArr, _ := appsAny.([]any)
	appMaps := make([]map[string]any, 0, len(appsArr))
	for _, it := range appsArr {
		appMaps = append(appMaps, extractFields(parsedData{"root": it}, prefixMappings(m)))
	}
	token, _ := asString(pathGet(parsed, tokenPath))

	page := pageMappings{Apps: []any{0, 0, 0}, Token: []any{0, 0, 7, 1}}
	more, err := checkFinished(ctx, c, opts.CallOptions, lang, country, num, appMaps, token, page)
	if err != nil {
		return nil, err
	}

	apps := make([]App, 0, len(more))
	for _, mm := range more {
		b, err := json.Marshal(mm)
		if err != nil {
			return nil, err
		}
		var a App
		if err := json.Unmarshal(b, &a); err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}

	if opts.FullDetail {
		out := make([]App, 0, len(apps))
		for _, a := range apps {
			full, err := c.App(ctx, AppOptions{CallOptions: opts.CallOptions, AppID: a.AppID, Lang: lang, Country: country})
			if err != nil {
				return nil, err
			}
			out = append(out, full)
		}
		return out, nil
	}

	c.cacheSet("developer", cacheOpts, apps)
	return apps, nil
}
