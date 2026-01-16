package gplay

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestIntegrationSmoke(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("set INTEGRATION=1 to run")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	c, err := NewClient(ClientOptions{Timeout: 25 * time.Second, RetryCount: 2, RetryWait: 700 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}

	call := CallOptions{Throttle: 1}

	app, err := c.App(ctx, AppOptions{CallOptions: call, AppID: "com.sgn.pandapop.gp", Lang: "en", Country: "us"})
	if err != nil {
		t.Fatalf("app: %v", err)
	}
	if app.AppID != "com.sgn.pandapop.gp" {
		t.Fatalf("unexpected appId: %q", app.AppID)
	}
	if app.Title == "" {
		t.Fatalf("empty title")
	}

	apps, err := c.Search(ctx, SearchOptions{CallOptions: call, Term: "netflix", Num: 10, Lang: "en", Country: "us"})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(apps) != 10 {
		t.Fatalf("search expected 10 results, got %d", len(apps))
	}

	if os.Getenv("INTEGRATION_QUALITY") == "1" {
		apps, err := c.Search(ctx, SearchOptions{CallOptions: call, Term: "gmail", Num: 20, Lang: "en", Country: "us"})
		if err != nil {
			t.Fatalf("search quality: %v", err)
		}
		found := false
		for _, a := range apps {
			if a.AppID == "com.google.android.gm" {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected gmail results to include com.google.android.gm")
		}
	}

	paginated, err := c.Search(ctx, SearchOptions{CallOptions: call, Term: "p", Num: 60, Lang: "en", Country: "us"})
	if err != nil {
		t.Fatalf("search pagination: %v", err)
	}
	if len(paginated) != 60 {
		t.Fatalf("search pagination expected 60 results, got %d", len(paginated))
	}

	devApps, err := c.Developer(ctx, DeveloperOptions{CallOptions: call, DevID: app.DeveloperID, Num: 10, Lang: "en", Country: "us"})
	if err != nil {
		t.Fatalf("developer: %v", err)
	}
	if len(devApps) == 0 {
		t.Fatalf("developer returned 0")
	}

	rev, err := c.Reviews(ctx, ReviewsOptions{CallOptions: call, AppID: "com.facebook.katana", Lang: "en", Country: "us", Paginate: true})
	if err != nil {
		t.Fatalf("reviews: %v", err)
	}
	if len(rev.Data) == 0 || rev.NextPaginationToken == nil {
		t.Fatalf("reviews paginate expected data+token")
	}

	cats, err := c.Categories(ctx, CategoriesOptions{CallOptions: call})
	if err != nil {
		t.Fatalf("categories: %v", err)
	}
	if len(cats) <= 1 {
		t.Fatalf("expected >1 category id, got %d", len(cats))
	}
}
