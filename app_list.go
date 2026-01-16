package gplay

import (
	"strconv"
	"strings"
)

func extractDeveloperID(link string) string {
	parts := strings.Split(link, "?id=")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

type appListExtractor struct{}

func (a appListExtractor) itemMappings() map[string]fieldSpec {
	return map[string]fieldSpec{
		"title": {Path: []any{2}},
		"appId": {Path: []any{12, 0}},
		"url": {Path: []any{9, 4, 2}, Fn: func(input any, _ parsedData) any {
			p, _ := asString(input)
			return resolveURL(BaseURL, p)
		}},
		"icon":      {Path: []any{1, 1, 0, 3, 2}},
		"developer": {Path: []any{4, 0, 0, 0}},
		"developerId": {Path: []any{4, 0, 0, 1, 4, 2}, Fn: func(input any, _ parsedData) any {
			s, _ := asString(input)
			return extractDeveloperID(s)
		}},
		"priceText": {Path: []any{7, 0, 3, 2, 1, 0, 2}, Fn: func(input any, _ parsedData) any {
			if input == nil {
				return "FREE"
			}
			s, _ := asString(input)
			return s
		}},
		"currency": {Path: []any{7, 0, 3, 2, 1, 0, 1}},
		"price": {Path: []any{7, 0, 3, 2, 1, 0, 2}, Fn: func(input any, _ parsedData) any {
			if input == nil {
				return float64(0)
			}
			s, _ := asString(input)
			if s == "" {
				return float64(0)
			}
			start := -1
			end := -1
			for i := 0; i < len(s); i++ {
				ch := s[i]
				if (ch >= '0' && ch <= '9') || ch == '.' || ch == ',' {
					if start == -1 {
						start = i
					}
					end = i + 1
				} else if start != -1 {
					break
				}
			}
			if start == -1 || end == -1 {
				return float64(0)
			}
			num := strings.ReplaceAll(s[start:end], ",", "")
			f, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return float64(0)
			}
			return f
		}},
		"free": {Path: []any{7, 0, 3, 2, 1, 0, 2}, Fn: func(input any, _ parsedData) any {
			return input == nil
		}},
		"summary":   {Path: []any{4, 1, 1, 1, 1}},
		"scoreText": {Path: []any{6, 0, 2, 1, 0}},
		"score":     {Path: []any{6, 0, 2, 1, 1}},
	}
}

func extractAppList(root []any, data any) []map[string]any {
	inputAny := pathGet(data, root)
	input, ok := inputAny.([]any)
	if !ok {
		return nil
	}
	m := appListExtractor{}.itemMappings()
	out := make([]map[string]any, 0, len(input))
	for _, it := range input {
		item := extractFields(parsedData{"root": it}, prefixMappings(m))
		out = append(out, item)
	}
	return out
}

func prefixMappings(m map[string]fieldSpec) map[string]fieldSpec {
	out := make(map[string]fieldSpec, len(m))
	for k, v := range m {
		p := make([]any, 0, 1+len(v.Path))
		p = append(p, "root")
		p = append(p, v.Path...)
		v.Path = p
		out[k] = v
	}
	return out
}
