package gplay

import "testing"

func TestPathGet(t *testing.T) {
	root := parsedData{
		"a": []any{map[string]any{"b": "c"}},
	}
	v := pathGet(root, []any{"a", 0, "b"})
	if s, _ := v.(string); s != "c" {
		t.Fatalf("expected c, got %#v", v)
	}
}

func TestParseListResponse(t *testing.T) {
	line3 := `[[null,null,"[1,2]"]]`
	body := []byte("x\ny\nz\n" + line3 + "\n")
	out, err := parseListResponse(body)
	if err != nil {
		t.Fatal(err)
	}
	arr, ok := out.([]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("expected [1,2], got %#v", out)
	}
}
