package gplay

import (
	"encoding/json"
	"errors"
)

func parseBatchedExecuteResponse(body []byte) (any, error) {
	if len(body) < 5 {
		return nil, errors.New("invalid batchexecute response")
	}
	trimmed := body
	if string(body[:5]) == ")]}'\n" {
		trimmed = body[5:]
	}
	var outer any
	if err := json.Unmarshal(trimmed, &outer); err != nil {
		return nil, err
	}
	return outer, nil
}

func parseBatchedInnerJSON(outer any) (any, error) {
	arr, ok := outer.([]any)
	if !ok || len(arr) == 0 {
		return nil, errors.New("unexpected batchexecute payload")
	}
	first, ok := arr[0].([]any)
	if !ok || len(first) < 3 {
		return nil, errors.New("unexpected batchexecute payload")
	}
	inner, ok := first[2].(string)
	if !ok {
		return nil, errors.New("unexpected batchexecute payload")
	}
	if inner == "null" {
		return nil, nil
	}
	var out any
	if err := json.Unmarshal([]byte(inner), &out); err != nil {
		return nil, err
	}
	return out, nil
}
