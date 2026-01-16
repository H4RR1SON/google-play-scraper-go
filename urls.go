package gplay

import (
	"net/url"
	"strings"
)

func resolveURL(base string, p string) string {
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	rel, err := url.Parse(p)
	if err != nil {
		return ""
	}
	return u.ResolveReference(rel).String()
}

func encodeValues(v url.Values) string {
	return strings.ReplaceAll(v.Encode(), "+", "%20")
}

func queryEscape(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}
