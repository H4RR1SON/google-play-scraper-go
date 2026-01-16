package gplay

import (
	"encoding/json"
	"sync"
	"time"
)

type memoCache struct {
	mu     sync.Mutex
	maxAge time.Duration
	max    int
	m      map[string]memoEntry
	order  []string
}

type memoEntry struct {
	expiresAt time.Time
	value     []byte
}

type MemoizeOptions struct {
	MaxAge time.Duration
	Max    int
}

func MemoizedClient(opts MemoizeOptions) *Client {
	c := MustNewClient(ClientOptions{Timeout: 15 * time.Second})
	maxAge := opts.MaxAge
	if maxAge == 0 {
		maxAge = 5 * time.Minute
	}
	max := opts.Max
	if max == 0 {
		max = 1000
	}
	c.cache = &memoCache{maxAge: maxAge, max: max, m: map[string]memoEntry{}}
	return c
}

func (c *Client) cacheKey(method string, opts any) (string, error) {
	b, err := json.Marshal(opts)
	if err != nil {
		return "", err
	}
	return method + ":" + string(b), nil
}

func (mc *memoCache) get(key string) ([]byte, bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	ent, ok := mc.m[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(ent.expiresAt) {
		delete(mc.m, key)
		return nil, false
	}
	return ent.value, true
}

func (mc *memoCache) set(key string, value []byte) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.m[key] = memoEntry{expiresAt: time.Now().Add(mc.maxAge), value: value}
	mc.order = append(mc.order, key)
	if mc.max > 0 && len(mc.m) > mc.max {
		mc.evict()
	}
}

func (mc *memoCache) evict() {
	for len(mc.order) > 0 && len(mc.m) > mc.max {
		k := mc.order[0]
		mc.order = mc.order[1:]
		delete(mc.m, k)
	}
}

func (c *Client) cacheGet(method string, opts any, out any) (bool, error) {
	if c == nil || c.cache == nil {
		return false, nil
	}
	key, err := c.cacheKey(method, opts)
	if err != nil {
		return false, err
	}
	b, ok := c.cache.get(key)
	if !ok {
		return false, nil
	}
	if err := json.Unmarshal(b, out); err != nil {
		return false, nil
	}
	return true, nil
}

func (c *Client) cacheSet(method string, opts any, val any) {
	if c == nil || c.cache == nil {
		return
	}
	key, err := c.cacheKey(method, opts)
	if err != nil {
		return
	}
	b, err := json.Marshal(val)
	if err != nil {
		return
	}
	c.cache.set(key, b)
}
