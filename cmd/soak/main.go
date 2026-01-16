package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	gplay "github.com/facundoolano/google-play-scraper-go"
)

type checkFn func(ctx context.Context, c *gplay.Client) error

type runner struct {
	c       *gplay.Client
	ok      int
	failed  int
	callOps gplay.CallOptions
}

func main() {
	timeout := 6 * time.Minute
	if v := os.Getenv("SOAK_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			timeout = d
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := gplay.NewClient(gplay.ClientOptions{Timeout: 25 * time.Second, RetryCount: 2, RetryWait: 700 * time.Millisecond})
	if err != nil {
		fatal(err)
	}

	throttle := 1
	if os.Getenv("SOAK_THROTTLE") != "" {
		fmt.Sscanf(os.Getenv("SOAK_THROTTLE"), "%d", &throttle)
		if throttle < 1 {
			throttle = 1
		}
	}

	r := &runner{c: client, callOps: gplay.CallOptions{Throttle: throttle}}

	checks := []struct {
		name string
		fn   checkFn
	}{
		{"app(en/us)", func(ctx context.Context, c *gplay.Client) error {
			app, err := c.App(ctx, gplay.AppOptions{CallOptions: r.callOps, AppID: "com.sgn.pandapop.gp", Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			return validateApp(app)
		}},
		{"app(es/es)", func(ctx context.Context, c *gplay.Client) error {
			app, err := c.App(ctx, gplay.AppOptions{CallOptions: r.callOps, AppID: "com.sgn.pandapop.gp", Lang: "es", Country: "es"})
			if err != nil {
				return err
			}
			return validateApp(app)
		}},
		{"search(netflix)", func(ctx context.Context, c *gplay.Client) error {
			apps, err := c.Search(ctx, gplay.SearchOptions{CallOptions: r.callOps, Term: "netflix", Num: 15, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(apps) != 15 {
				return fmt.Errorf("expected 15, got %d", len(apps))
			}
			for _, a := range apps {
				if err := validateLiteApp(a); err != nil {
					return err
				}
			}
			return nil
		}},
		{"search(gmail quality)", func(ctx context.Context, c *gplay.Client) error {
			apps, err := c.Search(ctx, gplay.SearchOptions{CallOptions: r.callOps, Term: "gmail", Num: 20, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(apps) == 0 {
				return errors.New("empty search")
			}
			if !containsAppID(apps, "com.google.android.gm") {
				return errors.New("expected results to include com.google.android.gm")
			}
			return nil
		}},
		{"search(spaces)", func(ctx context.Context, c *gplay.Client) error {
			apps, err := c.Search(ctx, gplay.SearchOptions{CallOptions: r.callOps, Term: "Panda vs Zombies", Num: 10, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(apps) != 10 {
				return fmt.Errorf("expected 10, got %d", len(apps))
			}
			return nil
		}},
		{"search(pagination)", func(ctx context.Context, c *gplay.Client) error {
			apps, err := c.Search(ctx, gplay.SearchOptions{CallOptions: r.callOps, Term: "p", Num: 80, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(apps) != 80 {
				return fmt.Errorf("expected 80, got %d", len(apps))
			}
			return nil
		}},
		{"search(empty)", func(ctx context.Context, c *gplay.Client) error {
			apps, err := c.Search(ctx, gplay.SearchOptions{CallOptions: r.callOps, Term: "asdasdyxcnmjysalsaflaslf", Num: 20, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(apps) != 0 {
				return fmt.Errorf("expected 0, got %d", len(apps))
			}
			return nil
		}},
		{"list(top_free)", func(ctx context.Context, c *gplay.Client) error {
			apps, err := c.List(ctx, gplay.ListOptions{CallOptions: r.callOps, Category: gplay.CategoryApplication, Collection: gplay.CollectionTopFree, Num: 20, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(apps) != 20 {
				return fmt.Errorf("expected 20, got %d", len(apps))
			}
			for _, a := range apps {
				if err := validateLiteApp(a); err != nil {
					return err
				}
			}
			return nil
		}},
		{"list(fullDetail)", func(ctx context.Context, c *gplay.Client) error {
			apps, err := c.List(ctx, gplay.ListOptions{CallOptions: r.callOps, Category: gplay.CategoryApplication, Collection: gplay.CollectionTopFree, Num: 5, Lang: "en", Country: "us", FullDetail: true})
			if err != nil {
				return err
			}
			if len(apps) != 5 {
				return fmt.Errorf("expected 5, got %d", len(apps))
			}
			for _, a := range apps {
				if err := validateApp(a); err != nil {
					return err
				}
			}
			return nil
		}},
		{"developer(devId)", func(ctx context.Context, c *gplay.Client) error {
			base, err := c.App(ctx, gplay.AppOptions{CallOptions: r.callOps, AppID: "com.sgn.pandapop.gp", Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if base.DeveloperID == "" {
				return errors.New("missing developerId")
			}
			apps, err := c.Developer(ctx, gplay.DeveloperOptions{CallOptions: r.callOps, DevID: base.DeveloperID, Num: 20, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(apps) == 0 {
				return errors.New("developer returned 0")
			}
			return nil
		}},
		{"suggest", func(ctx context.Context, c *gplay.Client) error {
			out, err := c.Suggest(ctx, gplay.SuggestOptions{CallOptions: r.callOps, Term: "panda", Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(out) == 0 {
				return errors.New("empty suggest")
			}
			return nil
		}},
		{"reviews(paginate)", func(ctx context.Context, c *gplay.Client) error {
			out, err := c.Reviews(ctx, gplay.ReviewsOptions{CallOptions: r.callOps, AppID: "com.facebook.katana", Lang: "en", Country: "us", Paginate: true})
			if err != nil {
				return err
			}
			if len(out.Data) == 0 || out.NextPaginationToken == nil {
				return errors.New("expected data + nextPaginationToken")
			}
			out2, err := c.Reviews(ctx, gplay.ReviewsOptions{CallOptions: r.callOps, AppID: "com.facebook.katana", Lang: "en", Country: "us", Paginate: true, NextPaginationToken: out.NextPaginationToken})
			if err != nil {
				return err
			}
			if len(out2.Data) == 0 {
				return errors.New("second page empty")
			}
			if out.Data[0].ID == out2.Data[0].ID {
				return errors.New("page2 appears identical to page1")
			}
			return nil
		}},
		{"similar", func(ctx context.Context, c *gplay.Client) error {
			apps, err := c.Similar(ctx, gplay.SimilarOptions{CallOptions: r.callOps, AppID: "com.sgn.pandapop.gp", Lang: "en", Country: "us", Num: 15})
			if err != nil {
				return err
			}
			if len(apps) != 15 {
				return fmt.Errorf("expected 15, got %d", len(apps))
			}
			return nil
		}},
		{"permissions", func(ctx context.Context, c *gplay.Client) error {
			out, err := c.Permissions(ctx, gplay.PermissionsOptions{CallOptions: r.callOps, AppID: "com.facebook.katana", Lang: "en", Country: "us", Short: false})
			if err != nil {
				return err
			}
			if len(out.Items) == 0 {
				return errors.New("permissions empty")
			}
			return nil
		}},
		{"datasafety", func(ctx context.Context, c *gplay.Client) error {
			out, err := c.DataSafety(ctx, gplay.DataSafetyOptions{CallOptions: r.callOps, AppID: "com.snapchat.android", Lang: "en"})
			if err != nil {
				return err
			}
			if len(out.CollectedData) == 0 {
				return errors.New("expected some collectedData")
			}
			return nil
		}},
		{"categories", func(ctx context.Context, c *gplay.Client) error {
			cats, err := c.Categories(ctx, gplay.CategoriesOptions{CallOptions: r.callOps})
			if err != nil {
				return err
			}
			if len(cats) < 5 {
				return fmt.Errorf("expected >=5 category ids, got %d", len(cats))
			}
			return nil
		}},
	}

	iters := 1
	if v := os.Getenv("SOAK_ITERS"); v != "" {
		fmt.Sscanf(v, "%d", &iters)
		if iters < 1 {
			iters = 1
		}
		if iters > 10 {
			iters = 10
		}
	}

	for i := 0; i < iters; i++ {
		for _, tc := range checks {
			r.step(ctx, tc.name, tc.fn)
		}
	}

	fmt.Printf("summary: ok=%d failed=%d\n", r.ok, r.failed)
	if r.failed > 0 {
		os.Exit(1)
	}
}

func (r *runner) step(ctx context.Context, name string, fn checkFn) {
	start := time.Now()
	err := fn(ctx, r.c)
	d := time.Since(start)
	if err != nil {
		r.failed++
		fmt.Printf("FAIL %-22s (%s): %s\n", name, d.Truncate(time.Millisecond), err.Error())
		return
	}
	r.ok++
	fmt.Printf("OK   %-22s (%s)\n", name, d.Truncate(time.Millisecond))
}

func containsAppID(apps []gplay.App, id string) bool {
	for _, a := range apps {
		if a.AppID == id {
			return true
		}
	}
	return false
}

func validateLiteApp(a gplay.App) error {
	if a.AppID == "" {
		return errors.New("empty appId")
	}
	if a.Title == "" {
		return errors.New("empty title")
	}
	if a.URL == "" || !strings.HasPrefix(a.URL, "https://") {
		return errors.New("invalid url")
	}
	if a.Icon == "" {
		return errors.New("empty icon")
	}
	return nil
}

func validateApp(a gplay.App) error {
	if err := validateLiteApp(a); err != nil {
		return err
	}
	if a.Summary == "" {
		return errors.New("empty summary")
	}
	if a.MinInstalls == nil {
		return errors.New("missing minInstalls")
	}
	if a.Reviews == nil {
		return errors.New("missing reviews")
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
