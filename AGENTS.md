# PROJECT KNOWLEDGE BASE

**Project:** google-play-scraper-go

## OVERVIEW
Go scraper for Google Play web pages/endpoints (port of facundoolano/google-play-scraper).

## STRUCTURE
```
./
├── *.go                 # library implementation (single-package)
├── *_test.go            # unit + opt-in integration tests
├── cmd/
│   ├── smoke/           # quick real-request sanity check
│   └── soak/            # broader real-request coverage
└── .github/workflows/   # CI (go test)
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| HTTP behavior (timeouts, proxy, retries) | client.go, request.go | Retry on 429/503; proxy via ClientOptions.ProxyURL |
| Throttling | throttle.go | Simple per-client throttle gate |
| Play HTML parsing (ds:* extraction) | scriptdata.go | Parses AF_initDataCallback + AF_dataServiceRequests |
| App details parsing | app.go, mappinghelpers.go | Field mapping + normalization |
| Search | search.go, pages.go, app_list.go | Pagination uses batchexecute qnKhOb |
| List/collections | list.go | Initial fetch via vyAe2 batchexecute |
| Reviews | reviews.go | Pagination token support |
| Permissions | permissions.go | Short vs full output |
| Data safety | datasafety.go | /store/apps/datasafety page parsing |
| Categories | categories.go | Anchor scan + regex fallback |
| Memoization | memoize.go | MemoizedClient (TTL + max entries) |

## CONVENTIONS (PROJECT-SPECIFIC)
- Public API is primarily `(*Client).Method(ctx, Options)`; convenience wrappers are `Fetch*` in api.go.
- Options embed `CallOptions` for per-call throttling and custom headers.
- Real-network tests are opt-in (env-gated). Default `go test ./...` must remain deterministic.

## ANTI-PATTERNS (THIS PROJECT)
- Do not rely on exact search ordering/IDs in tests (results change by geo/time).
- Do not run soak/smoke in CI by default.
- Keep fixes minimal when Play HTML/JSON shape changes; avoid refactors while patching scraping.

## COMMANDS
```bash
# Unit tests (deterministic)
go test ./...

go vet ./...

# Integration test (real Google Play requests)
INTEGRATION=1 go test -run TestIntegrationSmoke -v

# Optional “quality” assertions (may be more flaky)
INTEGRATION=1 INTEGRATION_QUALITY=1 go test -run TestIntegrationSmoke -v

# Local smoke tests (real requests)
go run ./cmd/smoke
SMOKE_EXTENDED=1 go run ./cmd/smoke

# Soak tests (real requests, multiple iterations)
go run ./cmd/soak
SOAK_ITERS=3 go run ./cmd/soak
SOAK_THROTTLE=1 SOAK_TIMEOUT=6m go run ./cmd/soak
```

## NOTES
- Scraping relies on undocumented endpoints and is expected to break occasionally.
- Expect 429/503 during bursts; keep `Throttle` conservative for stable behavior.
