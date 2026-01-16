# cmd/

## OVERVIEW
Executable-only utilities for real-request validation (not part of the library API).

## STRUCTURE
```
cmd/
├── smoke/   # quick, focused real-request checks (can be extended)
└── soak/    # broader coverage + output-shape validation, supports multiple iterations
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Quick sanity check before releases | cmd/smoke/main.go | Runs core methods once; `SMOKE_EXTENDED=1` adds heavier checks |
| Repeatability / multi-locale checks | cmd/soak/main.go | Runs a battery of calls; prints per-check timings |

## CONVENTIONS
- Both CLIs should default to conservative throttling.
- CLIs must exit non-zero on failures (usable in local scripts).

## KNOBS
| Tool | Env | Meaning |
|------|-----|---------|
| smoke | `SMOKE_EXTENDED=1` | Adds pagination/fullDetail checks |
| soak | `SOAK_ITERS` | Repeat the whole check suite (clamped to 1..10) |
| soak | `SOAK_THROTTLE` | Requests/sec throttle (min 1) |
| soak | `SOAK_TIMEOUT` | Overall timeout (e.g. `6m`, `10m`) |

## RUNNING
See root `./AGENTS.md` for the canonical command list.
