package gplay

import (
	"encoding/json"
	"regexp"

	"github.com/dop251/goja"
)

var (
	afInitDataCallbackScriptRe = regexp.MustCompile(`(?s)>AF_initDataCallback[\s\S]*?<\/script`)
	dsKeyRe                    = regexp.MustCompile(`(ds:.*?)'`)
	dsValueRe                  = regexp.MustCompile(`(?s)data:([\s\S]*?),\s*sideChannel:\s*\{\}\}\);<\/`)

	afDataServiceRequestsRe = regexp.MustCompile(`(?s); var AF_dataServiceRequests[\s\S]*?; var AF_initDataChunkQueue`)
	serviceRequestsValueRe  = regexp.MustCompile(`(?s)\{'ds:[\s\S]*\}\}`)
)

type parsedData map[string]any

func parseScriptData(html []byte) parsedData {
	m := parsedData{}

	matches := afInitDataCallbackScriptRe.FindAllSubmatch(html, -1)
	if len(matches) == 0 {
		m["serviceRequestData"] = map[string]any{}
		return m
	}

	for _, match := range matches {
		script := match[0]
		keyMatch := dsKeyRe.FindSubmatch(script)
		valMatch := dsValueRe.FindSubmatch(script)
		if len(keyMatch) != 2 || len(valMatch) != 2 {
			continue
		}
		key := string(keyMatch[1])
		var value any
		if err := json.Unmarshal(valMatch[1], &value); err != nil {
			continue
		}
		m[key] = value
	}

	m["serviceRequestData"] = parseServiceRequests(html)
	return m
}

func parseServiceRequests(html []byte) map[string]any {
	matches := afDataServiceRequestsRe.FindAllSubmatch(html, -1)
	if len(matches) == 0 {
		return map[string]any{}
	}
	data := matches[0][0]
	valueMatch := serviceRequestsValueRe.FindSubmatch(data)
	if len(valueMatch) != 1 {
		return map[string]any{}
	}
	literal := string(valueMatch[0])
	rt := goja.New()
	v, err := rt.RunString("(" + literal + ")")
	if err != nil {
		return map[string]any{}
	}
	out, ok := v.Export().(map[string]any)
	if !ok {
		return map[string]any{}
	}
	return out
}

type serviceRequestSpec struct {
	Path                []any
	UseServiceRequestID string
}

func extractDataWithServiceRequestID(data parsedData, spec serviceRequestSpec) any {
	srd, _ := data["serviceRequestData"].(map[string]any)
	if len(srd) == 0 {
		return pathGet(data, spec.Path)
	}

	var dsRoot string
	for k, v := range srd {
		vm, ok := v.(map[string]any)
		if !ok {
			continue
		}
		id, _ := vm["id"].(string)
		if id == spec.UseServiceRequestID {
			dsRoot = k
			break
		}
	}
	if dsRoot == "" {
		return pathGet(data, spec.Path)
	}

	fullPath := make([]any, 0, 1+len(spec.Path))
	fullPath = append(fullPath, dsRoot)
	fullPath = append(fullPath, spec.Path...)
	return pathGet(data, fullPath)
}
