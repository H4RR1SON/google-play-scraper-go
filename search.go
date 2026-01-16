package gplay

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
)

func (c *Client) Search(ctx context.Context, opts SearchOptions) ([]App, error) {
	if opts.Term == "" {
		return nil, errors.New("Search term missing")
	}
	if opts.Num > 0 && opts.Num > 250 {
		return nil, errors.New("The number of results can't exceed 250")
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
		num = 20
	}
	cacheOpts := opts
	cacheOpts.Lang = lang
	cacheOpts.Country = country
	cacheOpts.Num = num
	if c != nil && c.cache != nil {
		var cached []App
		hit, err := c.cacheGet("search", cacheOpts, &cached)
		if err != nil {
			return nil, err
		}
		if hit {
			return cached, nil
		}
	}

	price := 0
	switch opts.Price {
	case SearchPriceFree:
		price = 1
	case SearchPricePaid:
		price = 2
	default:
		price = 0
	}

	qs := url.Values{}
	qs.Set("q", opts.Term)
	qs.Set("hl", lang)
	qs.Set("gl", country)
	qs.Set("price", string(rune('0'+price)))
	pageURL := "/work/search?" + encodeValues(qs)
	body, _, err := c.do(ctx, requestOptions{URL: pageURL, Headers: opts.Headers}, opts.Throttle)
	if err != nil {
		return nil, err
	}

	parsed := parseScriptData(body)
	sectionsAny := pathGet(parsed, []any{"ds:1", 0, 1, 0, 0})
	sections, _ := sectionsAny.([]any)
	if len(sections) == 0 {
		return nil, nil
	}

	appsAny := pathGet(parsed, []any{"ds:1", 0, 1, 0, 0, 0})
	appsArr, _ := appsAny.([]any)

	token := ""
	for _, sec := range sections {
		secArr, ok := sec.([]any)
		if !ok {
			continue
		}
		if len(secArr) >= 2 {
			if s, ok := secArr[1].(string); ok {
				token = s
				break
			}
		}
	}

	m := map[string]fieldSpec{
		"title":       {Path: []any{2}},
		"appId":       {Path: []any{12, 0}},
		"url":         {Path: []any{9, 4, 2}, Fn: func(input any, _ parsedData) any { p, _ := asString(input); return resolveURL(BaseURL, p) }},
		"icon":        {Path: []any{1, 1, 0, 3, 2}},
		"developer":   {Path: []any{4, 0, 0, 0}},
		"developerId": {Path: []any{4, 0, 0, 1, 4, 2}, Fn: func(input any, _ parsedData) any { s, _ := asString(input); return extractDeveloperID(s) }},
		"currency":    {Path: []any{7, 0, 3, 2, 1, 0, 1}},
		"price":       {Path: []any{7, 0, 3, 2, 1, 0, 0}, Fn: func(input any, _ parsedData) any { v, _ := asFloat(input); return v / 1000000 }},
		"free":        {Path: []any{7, 0, 3, 2, 1, 0, 0}, Fn: func(input any, _ parsedData) any { v, _ := asFloat(input); return v == 0 }},
		"summary":     {Path: []any{4, 1, 1, 1, 1}},
		"scoreText":   {Path: []any{6, 0, 2, 1, 0}},
		"score":       {Path: []any{6, 0, 2, 1, 1}},
	}

	appMaps := make([]map[string]any, 0, len(appsArr))
	for _, it := range appsArr {
		fields := extractFields(parsedData{"root": it}, prefixMappings(m))
		appMaps = append(appMaps, fields)
	}

	page := pageMappings{Apps: []any{0, 0, 0}, Token: []any{0, 0, 7, 1}}
	more, err := checkFinished(ctx, c, opts.CallOptions, lang, country, num, appMaps, token, page)
	if err != nil {
		return nil, err
	}
	apps := make([]App, 0, len(more))
	for _, m := range more {
		b, err := json.Marshal(m)
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

	c.cacheSet("search", cacheOpts, apps)
	return apps, nil
}
