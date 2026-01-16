package gplay

import (
	"encoding/json"
	"testing"
)

func TestParseBatchedExecuteResponseAndInner(t *testing.T) {
	outerObj := []any{
		[]any{"wrb.fr", "1", "[1,2]"},
	}
	b, err := json.Marshal(outerObj)
	if err != nil {
		t.Fatal(err)
	}
	resp := append([]byte(")]}'\n"), b...)

	outer, err := parseBatchedExecuteResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	inner, err := parseBatchedInnerJSON(outer)
	if err != nil {
		t.Fatal(err)
	}
	arr, ok := inner.([]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("expected [1,2], got %#v", inner)
	}
}
