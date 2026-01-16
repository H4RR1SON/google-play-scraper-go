package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	gplay "github.com/facundoolano/google-play-scraper-go"
)

type runner struct {
	client   *gplay.Client
	callOpts gplay.CallOptions
	failures int
}

func main() {
	timeout := 2 * time.Minute
	if os.Getenv("SMOKE_EXTENDED") == "1" {
		timeout = 5 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := gplay.NewClient(gplay.ClientOptions{Timeout: 25 * time.Second, RetryCount: 2, RetryWait: 700 * time.Millisecond})
	if err != nil {
		fatal(err)
	}

	r := &runner{client: client, callOpts: gplay.CallOptions{Throttle: 1}}

	var app gplay.App
	r.step("app", func() error {
		out, err := r.client.App(ctx, gplay.AppOptions{CallOptions: r.callOpts, AppID: "com.sgn.pandapop.gp", Lang: "en", Country: "us"})
		if err != nil {
			return err
		}
		if out.AppID == "" || out.URL == "" {
			return errors.New("app output missing appId/url")
		}
		app = out
		fmt.Printf("app: appId=%s title=%s dev=%s devId=%s\n", out.AppID, out.Title, out.Developer, out.DeveloperID)
		return nil
	})

	r.step("search", func() error {
		apps, err := r.client.Search(ctx, gplay.SearchOptions{CallOptions: r.callOpts, Term: "netflix", Num: 10, Lang: "en", Country: "us"})
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			return errors.New("search returned 0 apps")
		}
		fmt.Printf("search: term=netflix count=%d first=%s\n", len(apps), apps[0].AppID)
		return nil
	})

	r.step("list", func() error {
		apps, err := r.client.List(ctx, gplay.ListOptions{CallOptions: r.callOpts, Category: gplay.CategoryApplication, Collection: gplay.CollectionTopFree, Num: 10, Lang: "en", Country: "us"})
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			return errors.New("list returned 0 apps")
		}
		fmt.Printf("list: TOP_FREE count=%d first=%s\n", len(apps), apps[0].AppID)
		return nil
	})

	r.step("developer(devId)", func() error {
		devID := app.DeveloperID
		if devID == "" {
			return errors.New("no devId available from app output")
		}
		apps, err := r.client.Developer(ctx, gplay.DeveloperOptions{CallOptions: r.callOpts, DevID: devID, Num: 10, Lang: "en", Country: "us"})
		if err != nil {
			return err
		}
		fmt.Printf("developer: devId=%s count=%d\n", devID, len(apps))
		return nil
	})

	r.step("developer(name)", func() error {
		if app.Developer == "" {
			return errors.New("no developer name available from app output")
		}
		apps, err := r.client.Developer(ctx, gplay.DeveloperOptions{CallOptions: r.callOpts, DevID: app.Developer, Num: 10, Lang: "en", Country: "us"})
		if err != nil {
			return err
		}
		fmt.Printf("developer: name=%q count=%d\n", app.Developer, len(apps))
		return nil
	})

	r.step("suggest", func() error {
		out, err := r.client.Suggest(ctx, gplay.SuggestOptions{CallOptions: r.callOpts, Term: "panda", Lang: "en", Country: "us"})
		if err != nil {
			return err
		}
		fmt.Printf("suggest: panda -> %v\n", out)
		return nil
	})

	r.step("reviews(helpfulness)", func() error {
		out, err := r.client.Reviews(ctx, gplay.ReviewsOptions{CallOptions: r.callOpts, AppID: "com.facebook.katana", Lang: "en", Country: "us", Sort: gplay.SortHelpfulness, Num: 20})
		if err != nil {
			return err
		}
		if len(out.Data) == 0 {
			return errors.New("reviews returned 0")
		}
		fmt.Printf("reviews: count=%d token=%t\n", len(out.Data), out.NextPaginationToken != nil)
		return nil
	})

	r.step("reviews(paginate)", func() error {
		page1, err := r.client.Reviews(ctx, gplay.ReviewsOptions{CallOptions: r.callOpts, AppID: "com.facebook.katana", Lang: "en", Country: "us", Paginate: true})
		if err != nil {
			return err
		}
		if len(page1.Data) == 0 || page1.NextPaginationToken == nil {
			return errors.New("reviews paginate returned empty or missing token")
		}
		page2, err := r.client.Reviews(ctx, gplay.ReviewsOptions{CallOptions: r.callOpts, AppID: "com.facebook.katana", Lang: "en", Country: "us", Paginate: true, NextPaginationToken: page1.NextPaginationToken})
		if err != nil {
			return err
		}
		fmt.Printf("reviews(paginate): page1=%d page2=%d token2=%t\n", len(page1.Data), len(page2.Data), page2.NextPaginationToken != nil)
		return nil
	})

	r.step("similar", func() error {
		apps, err := r.client.Similar(ctx, gplay.SimilarOptions{CallOptions: r.callOpts, AppID: "com.sgn.pandapop.gp", Lang: "en", Country: "us", Num: 10})
		if err != nil {
			return err
		}
		fmt.Printf("similar: count=%d\n", len(apps))
		return nil
	})

	r.step("permissions(short)", func() error {
		out, err := r.client.Permissions(ctx, gplay.PermissionsOptions{CallOptions: r.callOpts, AppID: "com.facebook.katana", Lang: "en", Country: "us", Short: true})
		if err != nil {
			return err
		}
		fmt.Printf("permissions(short): count=%d\n", len(out.Names))
		return nil
	})

	r.step("permissions(full)", func() error {
		out, err := r.client.Permissions(ctx, gplay.PermissionsOptions{CallOptions: r.callOpts, AppID: "com.facebook.katana", Lang: "en", Country: "us", Short: false})
		if err != nil {
			return err
		}
		fmt.Printf("permissions(full): count=%d\n", len(out.Items))
		return nil
	})

	r.step("datasafety", func() error {
		out, err := r.client.DataSafety(ctx, gplay.DataSafetyOptions{CallOptions: r.callOpts, AppID: "com.snapchat.android", Lang: "en"})
		if err != nil {
			return err
		}
		fmt.Printf("datasafety: shared=%d collected=%d security=%d\n", len(out.SharedData), len(out.CollectedData), len(out.SecurityPractices))
		return nil
	})

	r.step("categories", func() error {
		out, err := r.client.Categories(ctx, gplay.CategoriesOptions{CallOptions: r.callOpts})
		if err != nil {
			return err
		}
		fmt.Printf("categories: count=%d\n", len(out))
		return nil
	})

	if os.Getenv("SMOKE_EXTENDED") == "1" {
		r.step("search(spaces)", func() error {
			out, err := r.client.Search(ctx, gplay.SearchOptions{CallOptions: r.callOpts, Term: "Panda vs Zombies", Num: 5, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(out) == 0 {
				return errors.New("search with spaces returned 0")
			}
			fmt.Printf("search(spaces): count=%d first=%s\n", len(out), out[0].AppID)
			return nil
		})

		r.step("search(pagination)", func() error {
			out, err := r.client.Search(ctx, gplay.SearchOptions{CallOptions: r.callOpts, Term: "p", Num: 60, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(out) != 60 {
				return fmt.Errorf("expected 60 results, got %d", len(out))
			}
			return nil
		})

		r.step("search(empty)", func() error {
			out, err := r.client.Search(ctx, gplay.SearchOptions{CallOptions: r.callOpts, Term: "asdasdyxcnmjysalsaflaslf", Num: 20, Lang: "en", Country: "us"})
			if err != nil {
				return err
			}
			if len(out) != 0 {
				return fmt.Errorf("expected 0 results, got %d", len(out))
			}
			return nil
		})

		r.step("list(fullDetail)", func() error {
			out, err := r.client.List(ctx, gplay.ListOptions{CallOptions: r.callOpts, Category: gplay.CategoryApplication, Collection: gplay.CollectionTopFree, Num: 3, Lang: "en", Country: "us", FullDetail: true})
			if err != nil {
				return err
			}
			if len(out) != 3 {
				return fmt.Errorf("expected 3 results, got %d", len(out))
			}
			for i, a := range out {
				if a.AppID == "" || a.URL == "" || a.Title == "" {
					return fmt.Errorf("fullDetail app[%d] missing required fields", i)
				}
			}
			return nil
		})

		r.step("developer(fullDetail)", func() error {
			if app.DeveloperID == "" {
				return errors.New("no devId available")
			}
			out, err := r.client.Developer(ctx, gplay.DeveloperOptions{CallOptions: r.callOpts, DevID: app.DeveloperID, Num: 3, Lang: "en", Country: "us", FullDetail: true})
			if err != nil {
				return err
			}
			if len(out) != 3 {
				return fmt.Errorf("expected 3 results, got %d", len(out))
			}
			for i, a := range out {
				if a.AppID == "" || a.URL == "" || a.Title == "" {
					return fmt.Errorf("fullDetail app[%d] missing required fields", i)
				}
			}
			return nil
		})

		r.step("similar(fullDetail)", func() error {
			out, err := r.client.Similar(ctx, gplay.SimilarOptions{CallOptions: r.callOpts, AppID: "com.sgn.pandapop.gp", Lang: "en", Country: "us", Num: 3, FullDetail: true})
			if err != nil {
				return err
			}
			if len(out) != 3 {
				return fmt.Errorf("expected 3 results, got %d", len(out))
			}
			for i, a := range out {
				if a.AppID == "" || a.URL == "" || a.Title == "" {
					return fmt.Errorf("fullDetail app[%d] missing required fields", i)
				}
			}
			return nil
		})
	}

	if os.Getenv("SMOKE_PRINT_JSON") == "1" {
		b, _ := json.MarshalIndent(app, "", "  ")
		fmt.Println(string(b))
	}

	if r.failures > 0 {
		os.Exit(1)
	}
}

func (r *runner) step(name string, fn func() error) {
	start := time.Now()
	err := fn()
	d := time.Since(start)
	if err != nil {
		r.failures++
		fmt.Printf("FAIL %-20s (%s): %s\n", name, d.Truncate(time.Millisecond), err.Error())
		return
	}
	fmt.Printf("OK   %-20s (%s)\n", name, d.Truncate(time.Millisecond))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
