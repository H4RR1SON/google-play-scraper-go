package gplay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var permissionsURLTemplate = "/_/PlayStoreUi/data/batchexecute?rpcids=qnKhOb&f.sid=-697906427155521722&bl=boq_playuiserver_20190903.08_p0&hl=%s&gl=%s&authuser&soc-app=121&soc-platform=1&soc-device=1&_reqid=1065213"

func (c *Client) Permissions(ctx context.Context, opts PermissionsOptions) (PermissionsResult, error) {
	if opts.AppID == "" {
		return PermissionsResult{}, errors.New("appId missing")
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
		var cached PermissionsResult
		hit, err := c.cacheGet("permissions", cacheOpts, &cached)
		if err != nil {
			return PermissionsResult{}, err
		}
		if hit {
			return cached, nil
		}
	}

	u := fmt.Sprintf(permissionsURLTemplate, queryEscape(lang), queryEscape(country))
	body := fmt.Sprintf("f.req=%%5B%%5B%%5B%%22xdSrCf%%22%%2C%%22%%5B%%5Bnull%%2C%%5B%%5C%%22%s%%5C%%22%%2C7%%5D%%2C%%5B%%5D%%5D%%5D%%22%%2Cnull%%2C%%221%%22%%5D%%5D%%5D", opts.AppID)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	for k, vv := range opts.Headers {
		for _, v := range vv {
			headers.Add(k, v)
		}
	}

	respBody, _, err := c.do(ctx, requestOptions{Method: http.MethodPost, URL: u, Body: []byte(body), Headers: headers}, opts.Throttle)
	if err != nil {
		return PermissionsResult{}, err
	}
	outer, err := parseBatchedExecuteResponse(respBody)
	if err != nil {
		return PermissionsResult{}, err
	}
	inner, err := parseBatchedInnerJSON(outer)
	if err != nil {
		return PermissionsResult{}, err
	}
	if inner == nil {
		return PermissionsResult{Short: opts.Short}, nil
	}

	b, err := json.Marshal(inner)
	if err != nil {
		return PermissionsResult{}, err
	}
	var data any
	if err := json.Unmarshal(b, &data); err != nil {
		return PermissionsResult{}, err
	}

	rootArr, _ := data.([]any)
	if opts.Short {
		common, ok := pathGet(rootArr, []any{int(PermissionGroupCommon)}).([]any)
		if !ok {
			return PermissionsResult{Short: true}, nil
		}
		names := make([]string, 0)
		for _, p := range common {
			arr, ok := p.([]any)
			if !ok || len(arr) == 0 {
				continue
			}
			name, _ := asString(arr[0])
			if name != "" {
				names = append(names, name)
			}
		}
		res := PermissionsResult{Short: true, Names: names}
		c.cacheSet("permissions", cacheOpts, res)
		return res, nil
	}

	items := make([]PermissionItem, 0)
	for _, group := range []PermissionGroup{PermissionGroupCommon, PermissionGroupOther} {
		gArr, ok := pathGet(rootArr, []any{int(group)}).([]any)
		if !ok {
			continue
		}
		for _, section := range gArr {
			secArr, ok := section.([]any)
			if !ok {
				continue
			}
			typeLabel, _ := asString(pathGet(secArr, []any{0}))
			perms, ok := pathGet(secArr, []any{2}).([]any)
			if !ok {
				continue
			}
			for _, perm := range perms {
				permArr, ok := perm.([]any)
				if !ok {
					continue
				}
				name, _ := asString(pathGet(permArr, []any{1}))
				if name == "" {
					continue
				}
				items = append(items, PermissionItem{Permission: name, Type: typeLabel})
			}
		}
	}

	res := PermissionsResult{Short: false, Items: items}
	c.cacheSet("permissions", cacheOpts, res)
	return res, nil
}
