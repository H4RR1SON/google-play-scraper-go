package gplay

import "testing"

func TestParseScriptData_AFInitDataCallback(t *testing.T) {
	html := []byte(`
<html><head>
<script>AF_initDataCallback({key: 'ds:5', data:[1,2,3], sideChannel: {}});</script>
; var AF_dataServiceRequests = {'ds:0':{id:'ag2B9c'}}; var AF_initDataChunkQueue
</head></html>
`)

	pd := parseScriptData(html)
	if _, ok := pd["ds:5"]; !ok {
		t.Fatalf("expected ds:5 key")
	}

	srd, ok := pd["serviceRequestData"].(map[string]any)
	if !ok {
		t.Fatalf("expected serviceRequestData map")
	}
	root, ok := srd["ds:0"].(map[string]any)
	if !ok {
		t.Fatalf("expected ds:0 entry")
	}
	if id, _ := root["id"].(string); id != "ag2B9c" {
		t.Fatalf("expected id ag2B9c, got %q", id)
	}
}

func TestExtractDataWithServiceRequestID(t *testing.T) {
	pd := parsedData{
		"ds:1": map[string]any{"foo": "bar"},
		"serviceRequestData": map[string]any{
			"ds:1": map[string]any{"id": "ag2B9c"},
		},
	}

	v := extractDataWithServiceRequestID(pd, serviceRequestSpec{Path: []any{"foo"}, UseServiceRequestID: "ag2B9c"})
	if s, _ := v.(string); s != "bar" {
		t.Fatalf("expected bar, got %#v", v)
	}
}
