package gplay

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) Categories(ctx context.Context, opts CategoriesOptions) ([]string, error) {
	if c != nil && c.cache != nil {
		var cached []string
		hit, err := c.cacheGet("categories", opts, &cached)
		if err != nil {
			return nil, err
		}
		if hit {
			return cached, nil
		}
	}
	body, _, err := c.do(ctx, requestOptions{URL: "/store/apps", Headers: opts.Headers}, opts.Throttle)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	const prefix = "/store/apps/category/"
	categoryIDs := make([]string, 0)
	seen := map[string]struct{}{}
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}
		if !strings.HasPrefix(href, prefix) {
			return
		}
		if strings.Contains(href, "?age=") {
			return
		}
		id := strings.TrimPrefix(href, prefix)
		if i := strings.IndexByte(id, '?'); i >= 0 {
			id = id[:i]
		}
		if id == "" {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		categoryIDs = append(categoryIDs, id)
	})

	if len(categoryIDs) < 5 {
		re := regexp.MustCompile(`/store/apps/category/([A-Z0-9_]+)`)
		for _, m := range re.FindAllStringSubmatch(string(body), -1) {
			if len(m) != 2 {
				continue
			}
			id := m[1]
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			categoryIDs = append(categoryIDs, id)
		}
	}

	if _, ok := seen["APPLICATION"]; !ok {
		categoryIDs = append(categoryIDs, "APPLICATION")
	}
	if len(categoryIDs) == 0 {
		return nil, errors.New("no categories found")
	}
	c.cacheSet("categories", opts, categoryIDs)
	return categoryIDs, nil
}
