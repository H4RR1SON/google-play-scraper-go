package gplay

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func descriptionHTMLLocalized(search any) string {
	v := pathGet(search, []any{12, 0, 0, 1})
	if s, ok := asString(v); ok && s != "" {
		return s
	}
	v = pathGet(search, []any{72, 0, 1})
	s, _ := asString(v)
	return s
}

func descriptionText(description string) string {
	withBreaks := strings.ReplaceAll(description, "<br>", "\r\n")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<div>" + withBreaks + "</div>"))
	if err != nil {
		return ""
	}
	return doc.Find("div").Text()
}

func normalizeAndroidVersion(androidVersionText string) string {
	if androidVersionText == "" {
		return "VARY"
	}
	parts := strings.SplitN(androidVersionText, " ", 2)
	n := parts[0]
	for _, ch := range n {
		if (ch < '0' || ch > '9') && ch != '.' {
			return "VARY"
		}
	}
	if n == "" {
		return "VARY"
	}
	return n
}

func buildHistogram(container any) map[string]int64 {
	zero := map[string]int64{"1": 0, "2": 0, "3": 0, "4": 0, "5": 0}
	arr, ok := container.([]any)
	if !ok || len(arr) < 6 {
		return zero
	}
	for i := 1; i <= 5; i++ {
		sub, ok := arr[i].([]any)
		if !ok || len(sub) < 2 {
			continue
		}
		val, ok := asInt64(sub[1])
		if !ok {
			continue
		}
		zero[string(rune('0'+i))] = val
	}
	return zero
}

func extractComments(data parsedData) []string {
	for _, dsKey := range []string{"ds:8", "ds:9"} {
		author := pathGet(data, []any{dsKey, 0, 0, 1, 0})
		version := pathGet(data, []any{dsKey, 0, 0, 10})
		date := pathGet(data, []any{dsKey, 0, 0, 5, 0})
		if author == nil || version == nil || date == nil {
			continue
		}
		commentsAny := pathGet(data, []any{dsKey, 0})
		commentsArr, ok := commentsAny.([]any)
		if !ok {
			continue
		}
		out := make([]string, 0, 5)
		for _, c := range commentsArr {
			if len(out) >= 5 {
				break
			}
			text, _ := asString(pathGet(c, []any{4}))
			if text != "" {
				out = append(out, text)
			}
		}
		return out
	}
	return nil
}

func extractCategoriesRec(v any, categories *[]AppCategory) {
	arr, ok := v.([]any)
	if !ok || len(arr) == 0 {
		return
	}
	if len(arr) >= 4 {
		if name, ok := asString(arr[0]); ok {
			var idPtr *string
			if id, ok := asString(arr[2]); ok {
				idCopy := id
				idPtr = &idCopy
			}
			*categories = append(*categories, AppCategory{Name: name, ID: idPtr})
			return
		}
	}
	for _, sub := range arr {
		extractCategoriesRec(sub, categories)
	}
}

func extractCategories(v any) []AppCategory {
	var out []AppCategory
	extractCategoriesRec(v, &out)
	return out
}
