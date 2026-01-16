package gplay

import (
	"context"
	"fmt"
	"net/http"
)

var qnKhObURLTemplate = "/_/PlayStoreUi/data/batchexecute?rpcids=qnKhOb&f.sid=-697906427155521722&bl=boq_playuiserver_20190903.08_p0&hl=%s&gl=%s&authuser&soc-app=121&soc-platform=1&soc-device=1&_reqid=1065213"

func qnKhObBody(numberOfApps int, token string) string {
	if numberOfApps <= 0 {
		numberOfApps = 100
	}
	if token == "" {
		token = "%token%"
	}
	return fmt.Sprintf("f.req=%%5B%%5B%%5B%%22qnKhOb%%22%%2C%%22%%5B%%5Bnull%%2C%%5B%%5B10%%2C%%5B10%%2C%d%%5D%%5D%%2Ctrue%%2Cnull%%2C%%5B96%%2C27%%2C4%%2C8%%2C57%%2C30%%2C110%%2C79%%2C11%%2C16%%2C49%%2C1%%2C3%%2C9%%2C12%%2C104%%2C55%%2C56%%2C51%%2C10%%2C34%%2C77%%5D%%5D%%2Cnull%%2C%%5C%%22%s%%5C%%22%%5D%%5D%%22%%2Cnull%%2C%%22generic%%22%%5D%%5D%%5D", numberOfApps, token)
}

type pageMappings struct {
	Apps  []any
	Token []any
}

func checkFinished(ctx context.Context, c *Client, opts CallOptions, lang, country string, num int, saved []map[string]any, nextToken string, mappings pageMappings) ([]map[string]any, error) {
	if num <= 0 {
		return nil, nil
	}
	if len(saved) >= num || nextToken == "" {
		if len(saved) > num {
			return saved[:num], nil
		}
		return saved, nil
	}
	if c == nil {
		c = DefaultClient
	}

	u := fmt.Sprintf(qnKhObURLTemplate, queryEscape(lang), queryEscape(country))
	body := qnKhObBody(100, nextToken)

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
		if len(saved) > num {
			return saved[:num], nil
		}
		return saved, nil
	}
	return processPages(ctx, c, opts, lang, country, num, saved, inner, mappings)
}

func processPages(ctx context.Context, c *Client, opts CallOptions, lang, country string, num int, saved []map[string]any, data any, mappings pageMappings) ([]map[string]any, error) {
	apps := extractAppList(mappings.Apps, data)
	tokenVal, _ := asString(pathGet(data, mappings.Token))
	all := append(saved, apps...)
	return checkFinished(ctx, c, opts, lang, country, num, all, tokenVal, mappings)
}
