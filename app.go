package gplay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) App(ctx context.Context, opts AppOptions) (App, error) {
	if opts.AppID == "" {
		return App{}, errors.New("appId missing")
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
		var cached App
		hit, err := c.cacheGet("app", cacheOpts, &cached)
		if err != nil {
			return App{}, err
		}
		if hit {
			return cached, nil
		}
	}

	qs := url.Values{}
	qs.Set("id", opts.AppID)
	qs.Set("hl", lang)
	qs.Set("gl", country)
	pageURL := "/store/apps/details?" + encodeValues(qs)
	body, _, err := c.do(ctx, requestOptions{URL: pageURL, Headers: opts.Headers}, opts.Throttle)
	if err != nil {
		return App{}, err
	}

	parsed := parseScriptData(body)

	mappings := map[string]fieldSpec{
		"title": {Path: []any{"ds:5", 1, 2, 0, 0}},
		"description": {Path: []any{"ds:5", 1, 2}, Fn: func(input any, _ parsedData) any {
			html := descriptionHTMLLocalized(input)
			text := descriptionText(html)
			return text
		}},
		"descriptionHTML": {Path: []any{"ds:5", 1, 2}, Fn: func(input any, _ parsedData) any { return descriptionHTMLLocalized(input) }},
		"summary":         {Path: []any{"ds:5", 1, 2, 73, 0, 1}},
		"installs":        {Path: []any{"ds:5", 1, 2, 13, 0}},
		"minInstalls":     {Path: []any{"ds:5", 1, 2, 13, 1}, Fn: func(input any, _ parsedData) any { v, _ := asInt64(input); return v }},
		"maxInstalls":     {Path: []any{"ds:5", 1, 2, 13, 2}, Fn: func(input any, _ parsedData) any { v, _ := asInt64(input); return v }},
		"score":           {Path: []any{"ds:5", 1, 2, 51, 0, 1}},
		"scoreText":       {Path: []any{"ds:5", 1, 2, 51, 0, 0}},
		"ratings":         {Path: []any{"ds:5", 1, 2, 51, 2, 1}, Fn: func(input any, _ parsedData) any { v, _ := asInt64(input); return v }},
		"reviews":         {Path: []any{"ds:5", 1, 2, 51, 3, 1}, Fn: func(input any, _ parsedData) any { v, _ := asInt64(input); return v }},
		"histogram":       {Path: []any{"ds:5", 1, 2, 51, 1}, Fn: func(input any, _ parsedData) any { return buildHistogram(input) }},
		"price": {Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 0, 0}, Fn: func(input any, _ parsedData) any {
			v, ok := asFloat(input)
			if !ok {
				return float64(0)
			}
			return v / 1000000
		}},
		"originalPrice": {Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 1, 0}, Fn: func(input any, _ parsedData) any {
			v, ok := asFloat(input)
			if !ok || v == 0 {
				return nil
			}
			return v / 1000000
		}},
		"discountEndDate": {Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 14, 1}},
		"free": {Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 0, 0}, Fn: func(input any, _ parsedData) any {
			v, ok := asFloat(input)
			return ok && v == 0
		}},
		"currency": {Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 0, 1}},
		"priceText": {Path: []any{"ds:5", 1, 2, 57, 0, 0, 0, 0, 1, 0, 2}, Fn: func(input any, _ parsedData) any {
			s, _ := asString(input)
			if s == "" {
				return "Free"
			}
			return s
		}},
		"available": {Path: []any{"ds:5", 1, 2, 18, 0}, Fn: func(input any, _ parsedData) any { return truthy(input) }},
		"offersIAP": {Path: []any{"ds:5", 1, 2, 19, 0}, Fn: func(input any, _ parsedData) any { return truthy(input) }},
		"IAPRange":  {Path: []any{"ds:5", 1, 2, 19, 0}},
		"androidVersion": {Path: []any{"ds:5", 1, 2, 140, 1, 1, 0, 0, 1}, FallbackPath: []any{"ds:5", 1, 2, -1, "141", 1, 1, 0, 0, 1}, Fn: func(input any, _ parsedData) any {
			s, _ := asString(input)
			return normalizeAndroidVersion(s)
		}},
		"androidVersionText": {Path: []any{"ds:5", 1, 2, 140, 1, 1, 0, 0, 1}, FallbackPath: []any{"ds:5", 1, 2, -1, "141", 1, 1, 0, 0, 1}, Fn: func(input any, _ parsedData) any {
			s, _ := asString(input)
			if s == "" {
				return "Varies with device"
			}
			return s
		}},
		"androidMaxVersion": {Path: []any{"ds:5", 1, 2, 140, 1, 1, 0, 1, 1}, FallbackPath: []any{"ds:5", 1, 2, -1, "141", 1, 1, 0, 1, 1}, Fn: func(input any, _ parsedData) any {
			s, _ := asString(input)
			return normalizeAndroidVersion(s)
		}},
		"developer": {Path: []any{"ds:5", 1, 2, 68, 0}},
		"developerId": {Path: []any{"ds:5", 1, 2, 68, 1, 4, 2}, Fn: func(input any, _ parsedData) any {
			s, _ := asString(input)
			parts := strings.Split(s, "id=")
			if len(parts) < 2 {
				return ""
			}
			return parts[1]
		}},
		"developerEmail":      {Path: []any{"ds:5", 1, 2, 69, 1, 0}},
		"developerWebsite":    {Path: []any{"ds:5", 1, 2, 69, 0, 5, 2}},
		"developerAddress":    {Path: []any{"ds:5", 1, 2, 69, 2, 0}},
		"developerLegalName":  {Path: []any{"ds:5", 1, 2, 69, 4, 0}},
		"developerLegalEmail": {Path: []any{"ds:5", 1, 2, 69, 4, 1, 0}},
		"developerLegalAddress": {Path: []any{"ds:5", 1, 2, 69}, Fn: func(input any, _ parsedData) any {
			v := pathGet(input, []any{4, 2, 0})
			s, _ := asString(v)
			if s == "" {
				return nil
			}
			return strings.ReplaceAll(s, "\n", ", ")
		}},
		"developerLegalPhoneNumber": {Path: []any{"ds:5", 1, 2, 69, 4, 3}},
		"privacyPolicy":             {Path: []any{"ds:5", 1, 2, 99, 0, 5, 2}},
		"developerInternalID": {Path: []any{"ds:5", 1, 2, 68, 1, 4, 2}, Fn: func(input any, _ parsedData) any {
			s, _ := asString(input)
			parts := strings.Split(s, "id=")
			if len(parts) < 2 {
				return ""
			}
			return parts[1]
		}},
		"genre":   {Path: []any{"ds:5", 1, 2, 79, 0, 0, 0}},
		"genreId": {Path: []any{"ds:5", 1, 2, 79, 0, 0, 2}},
		"categories": {Path: []any{"ds:5", 1, 2}, Fn: func(input any, _ parsedData) any {
			cats := extractCategories(pathGet(input, []any{118}))
			if len(cats) == 0 {
				name, _ := asString(pathGet(input, []any{79, 0, 0, 0}))
				id, _ := asString(pathGet(input, []any{79, 0, 0, 2}))
				idCopy := id
				cats = append(cats, AppCategory{Name: name, ID: &idCopy})
			}
			return cats
		}},
		"icon":        {Path: []any{"ds:5", 1, 2, 95, 0, 3, 2}},
		"headerImage": {Path: []any{"ds:5", 1, 2, 96, 0, 3, 2}},
		"screenshots": {Path: []any{"ds:5", 1, 2, 78, 0}, Fn: func(input any, _ parsedData) any {
			arr, ok := input.([]any)
			if !ok {
				return []string{}
			}
			out := make([]string, 0, len(arr))
			for _, it := range arr {
				u, _ := asString(pathGet(it, []any{3, 2}))
				if u != "" {
					out = append(out, u)
				}
			}
			return out
		}},
		"video":                    {Path: []any{"ds:5", 1, 2, 100, 0, 0, 3, 2}},
		"videoImage":               {Path: []any{"ds:5", 1, 2, 100, 1, 0, 3, 2}},
		"previewVideo":             {Path: []any{"ds:5", 1, 2, 100, 1, 2, 0, 2}},
		"contentRating":            {Path: []any{"ds:5", 1, 2, 9, 0}},
		"contentRatingDescription": {Path: []any{"ds:5", 1, 2, 9, 2, 1}},
		"adSupported":              {Path: []any{"ds:5", 1, 2, 48}, Fn: func(input any, _ parsedData) any { return truthy(input) }},
		"released":                 {Path: []any{"ds:5", 1, 2, 10, 0}},
		"updated": {Path: []any{"ds:5", 1, 2, 145, 0, 1, 0}, FallbackPath: []any{"ds:5", 1, 2, -1, "146", 0, 1, 0}, Fn: func(input any, _ parsedData) any {
			v, ok := asInt64(input)
			if !ok {
				return int64(0)
			}
			return v * 1000
		}},
		"version": {Path: []any{"ds:5", 1, 2, 140, 0, 0, 0}, FallbackPath: []any{"ds:5", 1, 2, -1, "141", 0, 0, 0}, Fn: func(input any, _ parsedData) any {
			s, _ := asString(input)
			if s == "" {
				return "VARY"
			}
			return s
		}},
		"recentChanges": {Path: []any{"ds:5", 1, 2, 144, 1, 1}, FallbackPath: []any{"ds:5", 1, 2, -1, "145", 1, 1}},
		"comments":      {Fn: func(_ any, data parsedData) any { return extractComments(data) }},
		"preregister": {Path: []any{"ds:5", 1, 2, 18, 0}, Fn: func(input any, _ parsedData) any {
			v, _ := asFloat(input)
			return v == 1
		}},
		"earlyAccessEnabled": {Path: []any{"ds:5", 1, 2, 18, 2}, Fn: func(input any, _ parsedData) any {
			_, ok := asString(input)
			return ok
		}},
		"isAvailableInPlayPass": {Path: []any{"ds:5", 1, 2, 62}, Fn: func(input any, _ parsedData) any { return truthy(input) }},
	}

	fields := extractFields(parsed, mappings)
	fields["appId"] = opts.AppID
	fields["url"] = c.withBaseURL(pageURL)

	b, err := json.Marshal(fields)
	if err != nil {
		return App{}, err
	}
	var app App
	if err := json.Unmarshal(b, &app); err != nil {
		return App{}, fmt.Errorf("unmarshal app: %w", err)
	}
	c.cacheSet("app", cacheOpts, app)
	return app, nil
}

func truthy(v any) bool {
	if v == nil {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case float64:
		return t != 0
	case string:
		return t != ""
	default:
		return true
	}
}
