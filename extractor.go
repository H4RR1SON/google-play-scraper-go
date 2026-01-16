package gplay

import (
	"fmt"
)

type fieldSpec struct {
	Path                []any
	FallbackPath        []any
	UseServiceRequestID string
	Fn                  func(input any, data parsedData) any
}

func extractFields(data parsedData, mappings map[string]fieldSpec) map[string]any {
	out := make(map[string]any, len(mappings))
	for k, spec := range mappings {
		var input any
		if spec.UseServiceRequestID != "" {
			input = extractDataWithServiceRequestID(data, serviceRequestSpec{Path: spec.Path, UseServiceRequestID: spec.UseServiceRequestID})
		} else {
			input = pathGet(data, spec.Path)
			if (input == nil) && len(spec.FallbackPath) > 0 {
				input = pathGet(data, spec.FallbackPath)
			}
		}
		if spec.Fn != nil {
			out[k] = spec.Fn(input, data)
		} else {
			out[k] = input
		}
	}
	return out
}

func pathGet(root any, path []any) any {
	if len(path) == 0 {
		return root
	}
	cur := root
	for _, p := range path {
		switch c := cur.(type) {
		case parsedData:
			cur = c[fmt.Sprint(p)]
		case map[string]any:
			cur = c[fmt.Sprint(p)]
		case []any:
			idx, ok := p.(int)
			if !ok {
				return nil
			}
			if idx < 0 || idx >= len(c) {
				return nil
			}
			cur = c[idx]
		default:
			return nil
		}
		if cur == nil {
			return nil
		}
	}
	return cur
}

func asString(v any) (string, bool) {
	s, ok := v.(string)
	return s, ok
}

func asFloat(v any) (float64, bool) {
	sw, ok := v.(float64)
	return sw, ok
}

func asInt64(v any) (int64, bool) {
	sw, ok := v.(float64)
	if ok {
		return int64(sw), true
	}
	iv, ok := v.(int64)
	if ok {
		return iv, true
	}
	return 0, false
}
