package gplay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (c *Client) Reviews(ctx context.Context, opts ReviewsOptions) (ReviewsResult, error) {
	if opts.AppID == "" {
		return ReviewsResult{}, errors.New("appId missing")
	}
	lang := opts.Lang
	if lang == "" {
		lang = "en"
	}
	country := opts.Country
	if country == "" {
		country = "us"
	}
	sort := opts.Sort
	if sort == 0 {
		sort = SortNewest
	}
	num := opts.Num
	if num == 0 {
		num = 150
	}
	cacheOpts := opts
	cacheOpts.Lang = lang
	cacheOpts.Country = country
	cacheOpts.Sort = sort
	cacheOpts.Num = num
	if c != nil && c.cache != nil {
		var cached ReviewsResult
		hit, err := c.cacheGet("reviews", cacheOpts, &cached)
		if err != nil {
			return ReviewsResult{}, err
		}
		if hit {
			return cached, nil
		}
	}

	requestType := "initial"
	token := "%token%"
	if opts.NextPaginationToken != nil {
		requestType = "paginated"
		token = *opts.NextPaginationToken
	}

	res, err := c.makeReviewsRequest(ctx, opts.CallOptions, lang, country, opts.AppID, sort, num, opts.Paginate, token, requestType, nil)
	if err != nil {
		return ReviewsResult{}, err
	}
	c.cacheSet("reviews", cacheOpts, res)
	return res, nil
}

func (c *Client) makeReviewsRequest(ctx context.Context, callOpts CallOptions, lang, country, appID string, sort Sort, num int, paginate bool, token string, requestType string, saved []Review) (ReviewsResult, error) {
	urlStr := fmt.Sprintf(qnKhObURLTemplate, queryEscape(lang), queryEscape(country))
	body := reviewsBody(appID, int(sort), 150, token, requestType)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	for k, vv := range callOpts.Headers {
		for _, v := range vv {
			headers.Add(k, v)
		}
	}

	respBody, _, err := c.do(ctx, requestOptions{Method: http.MethodPost, URL: urlStr, Body: []byte(body), Headers: headers}, callOpts.Throttle)
	if err != nil {
		return ReviewsResult{}, err
	}
	outer, err := parseBatchedExecuteResponse(respBody)
	if err != nil {
		return ReviewsResult{}, err
	}
	inner, err := parseBatchedInnerJSON(outer)
	if err != nil {
		return ReviewsResult{}, err
	}
	if inner == nil {
		return formatReviews(saved, num, nil), nil
	}

	b, err := json.Marshal(inner)
	if err != nil {
		return ReviewsResult{}, err
	}
	var data any
	if err := json.Unmarshal(b, &data); err != nil {
		return ReviewsResult{}, err
	}

	return c.processReviewsAndNext(ctx, callOpts, lang, country, appID, sort, num, paginate, data, saved)
}

func (c *Client) processReviewsAndNext(ctx context.Context, callOpts CallOptions, lang, country, appID string, sort Sort, num int, paginate bool, payload any, saved []Review) (ReviewsResult, error) {
	arr, ok := payload.([]any)
	if !ok || len(arr) == 0 {
		return formatReviews(saved, num, nil), nil
	}

	reviews := extractReviews(payload, appID)
	token, _ := asString(pathGet(payload, []any{1, 1}))
	acc := append(saved, reviews...)

	if !paginate && token != "" && len(acc) < num {
		return c.makeReviewsRequest(ctx, callOpts, lang, country, appID, sort, num, paginate, token, "paginated", acc)
	}

	var next *string
	if token != "" {
		t := token
		next = &t
	}
	return formatReviews(acc, num, next), nil
}

func formatReviews(reviews []Review, num int, token *string) ReviewsResult {
	out := reviews
	if len(out) > num {
		out = out[:num]
	}
	return ReviewsResult{Data: out, NextPaginationToken: token}
}

func reviewsBody(appID string, sort int, perRequest int, token string, requestType string) string {
	if requestType == "paginated" {
		return fmt.Sprintf("f.req=%%5B%%5B%%5B%%22UsvDTd%%22%%2C%%22%%5Bnull%%2Cnull%%2C%%5B2%%2C%d%%2C%%5B%d%%2Cnull%%2C%%5C%%22%s%%5C%%22%%5D%%2Cnull%%2C%%5B%%5D%%5D%%2C%%5B%%5C%%22%s%%5C%%22%%2C7%%5D%%5D%%22%%2Cnull%%2C%%22generic%%22%%5D%%5D%%5D", sort, perRequest, token, appID)
	}
	return fmt.Sprintf("f.req=%%5B%%5B%%5B%%22UsvDTd%%22%%2C%%22%%5Bnull%%2Cnull%%2C%%5B2%%2C%d%%2C%%5B%d%%2Cnull%%2Cnull%%5D%%2Cnull%%2C%%5B%%5D%%5D%%2C%%5B%%5C%%22%s%%5C%%22%%2C7%%5D%%5D%%22%%2Cnull%%2C%%22generic%%22%%5D%%5D%%5D", sort, perRequest, appID)
}

func extractReviews(payload any, appID string) []Review {
	root, ok := pathGet(payload, []any{0}).([]any)
	if !ok {
		return nil
	}
	out := make([]Review, 0, len(root))
	for _, it := range root {
		m := map[string]fieldSpec{
			"id":        {Path: []any{0}},
			"userName":  {Path: []any{1, 0}},
			"userImage": {Path: []any{1, 1, 3, 2}},
			"date":      {Path: []any{5}, Fn: func(input any, _ parsedData) any { return generateDate(input) }},
			"score":     {Path: []any{2}, Fn: func(input any, _ parsedData) any { v, _ := asInt64(input); return v }},
			"scoreText": {Path: []any{2}, Fn: func(input any, _ parsedData) any { v, _ := asInt64(input); return fmt.Sprint(v) }},
			"url": {Path: []any{0}, Fn: func(input any, _ parsedData) any {
				rid, _ := asString(input)
				return BaseURL + "/store/apps/details?id=" + appID + "&reviewId=" + rid
			}},
			"title":     {Fn: func(_ any, _ parsedData) any { return nil }},
			"text":      {Path: []any{4}},
			"replyDate": {Path: []any{7, 2}, Fn: func(input any, _ parsedData) any { return generateDate(input) }},
			"replyText": {Path: []any{7, 1}, Fn: func(input any, _ parsedData) any {
				s, _ := asString(input)
				if s == "" {
					return nil
				}
				return s
			}},
			"version": {Path: []any{10}, Fn: func(input any, _ parsedData) any {
				s, _ := asString(input)
				if s == "" {
					return nil
				}
				return s
			}},
			"thumbsUp":  {Path: []any{6}, Fn: func(input any, _ parsedData) any { v, _ := asInt64(input); return v }},
			"criterias": {Path: []any{12, 0}, Fn: func(input any, _ parsedData) any { return buildCriterias(input) }},
		}
		fields := extractFields(parsedData{"root": it}, prefixMappings(m))
		b, err := json.Marshal(fields)
		if err != nil {
			continue
		}
		var r Review
		if err := json.Unmarshal(b, &r); err != nil {
			continue
		}
		out = append(out, r)
	}
	return out
}

func buildCriterias(v any) []ReviewCriteria {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]ReviewCriteria, 0, len(arr))
	for _, it := range arr {
		criteria, _ := asString(pathGet(it, []any{0}))
		ratingVal := pathGet(it, []any{1, 0})
		var rating *int64
		if r, ok := asInt64(ratingVal); ok {
			rc := r
			rating = &rc
		}
		out = append(out, ReviewCriteria{Criteria: criteria, Rating: rating})
	}
	return out
}

func generateDate(v any) *string {
	arr, ok := v.([]any)
	if !ok || len(arr) == 0 {
		return nil
	}
	ms, ok := asInt64(arr[0])
	if !ok {
		return nil
	}
	last := "000"
	if len(arr) > 1 {
		if s, ok := asString(arr[1]); ok && s != "" {
			last = s
		} else if i, ok := asInt64(arr[1]); ok {
			last = fmt.Sprint(i)
		}
	}
	if len(last) > 3 {
		last = last[:3]
	}
	for len(last) < 3 {
		last += "0"
	}
	full := fmt.Sprintf("%d%s", ms, last)
	val, err := strconv.ParseInt(full, 10, 64)
	if err != nil {
		return nil
	}
	tm := time.UnixMilli(val).UTC().Format(time.RFC3339Nano)
	return &tm
}
