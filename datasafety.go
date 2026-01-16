package gplay

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
)

func (c *Client) DataSafety(ctx context.Context, opts DataSafetyOptions) (DataSafetyResult, error) {
	if opts.AppID == "" {
		return DataSafetyResult{}, errors.New("appId missing")
	}
	lang := opts.Lang
	if lang == "" {
		lang = "en"
	}
	cacheOpts := opts
	cacheOpts.Lang = lang
	if c != nil && c.cache != nil {
		var cached DataSafetyResult
		hit, err := c.cacheGet("datasafety", cacheOpts, &cached)
		if err != nil {
			return DataSafetyResult{}, err
		}
		if hit {
			return cached, nil
		}
	}

	qs := url.Values{}
	qs.Set("id", opts.AppID)
	qs.Set("hl", lang)
	pageURL := "/store/apps/datasafety?" + encodeValues(qs)

	body, _, err := c.do(ctx, requestOptions{URL: pageURL, Headers: opts.Headers}, opts.Throttle)
	if err != nil {
		return DataSafetyResult{}, err
	}
	parsed := parseScriptData(body)

	mappings := map[string]fieldSpec{
		"dataShared":        {Path: []any{"ds:3", 1, 2, 1, 138, 4, 0, 0}, Fn: func(input any, _ parsedData) any { return mapDataEntries(input) }},
		"dataCollected":     {Path: []any{"ds:3", 1, 2, 1, 138, 4, 1, 0}, Fn: func(input any, _ parsedData) any { return mapDataEntries(input) }},
		"securityPractices": {Path: []any{"ds:3", 1, 2, 1, 138, 9, 2}, Fn: func(input any, _ parsedData) any { return mapSecurityPractices(input) }},
		"privacyPolicyUrl":  {Path: []any{"ds:3", 1, 2, 1, 100, 0, 5, 2}},
	}
	fields := extractFields(parsed, mappings)

	b, err := json.Marshal(fields)
	if err != nil {
		return DataSafetyResult{}, err
	}
	var out DataSafetyResult
	if err := json.Unmarshal(b, &out); err != nil {
		return DataSafetyResult{}, err
	}
	c.cacheSet("datasafety", cacheOpts, out)
	return out, nil
}

func mapSecurityPractices(v any) []SecurityPractice {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]SecurityPractice, 0, len(arr))
	for _, it := range arr {
		practice, _ := asString(pathGet(it, []any{1}))
		desc, _ := asString(pathGet(it, []any{2, 1}))
		if practice == "" && desc == "" {
			continue
		}
		out = append(out, SecurityPractice{Practice: practice, Description: desc})
	}
	return out
}

func mapDataEntries(v any) []DataSafetyEntry {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]DataSafetyEntry, 0)
	for _, data := range arr {
		typeStr, _ := asString(pathGet(data, []any{0, 1}))
		details, ok := pathGet(data, []any{4}).([]any)
		if !ok {
			continue
		}
		for _, detail := range details {
			dataStr, _ := asString(pathGet(detail, []any{0}))
			optional := truthy(pathGet(detail, []any{1}))
			purpose, _ := asString(pathGet(detail, []any{2}))
			out = append(out, DataSafetyEntry{Data: dataStr, Optional: optional, Purpose: purpose, Type: typeStr})
		}
	}
	return out
}
