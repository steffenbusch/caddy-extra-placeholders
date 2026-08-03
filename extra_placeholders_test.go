package extraplaceholders

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/shirou/gopsutil/v4/load"
)

func TestValidate(t *testing.T) {
	t.Run("valid range", func(t *testing.T) {
		e := ExtraPlaceholders{RandIntMin: 10, RandIntMax: 20}
		if err := e.Validate(); err != nil {
			t.Fatalf("Validate() returned error: %v", err)
		}
	})

	t.Run("equal bounds", func(t *testing.T) {
		e := ExtraPlaceholders{RandIntMin: 10, RandIntMax: 10}
		if err := e.Validate(); err == nil {
			t.Fatal("Validate() returned nil error, want error")
		}
	})

	t.Run("max below min", func(t *testing.T) {
		e := ExtraPlaceholders{RandIntMin: 20, RandIntMax: 10}
		if err := e.Validate(); err == nil {
			t.Fatal("Validate() returned nil error, want error")
		}
	})
}

func TestSetHTTPRequestURLPlaceholders(t *testing.T) {
	repl := caddy.NewReplacer()
	req := httptest.NewRequest(http.MethodGet, "/search?q=a+b&x=1", nil)
	req.Host = "example.com"
	req.TLS = &tls.ConnectionState{}

	ExtraPlaceholders{}.setHTTPRequestURLPlaceholders(repl, req)

	const want = "https%3A%2F%2Fexample.com%2Fsearch%3Fq%3Da%2Bb%26x%3D1"
	if got := mustGetString(t, repl, "extra.http.request.url.query_escaped"); got != want {
		t.Fatalf("query_escaped URL = %q, want %q", got, want)
	}
}

func TestSetTimePlaceholders(t *testing.T) {
	repl := caddy.NewReplacer()
	repl.Set("test.format", "2006/01/02 15:04")

	fixed := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.FixedZone("CET", 3600))
	isoYear, isoWeek := fixed.ISOWeek()

	e := ExtraPlaceholders{TimeFormatCustom: "{test.format}"}
	e.setTimePlaceholders(repl, fixed, false)

	assertReplacerValue(t, repl, "extra.time.now.month", 1)
	assertReplacerValue(t, repl, "extra.time.now.month_padded", "01")
	assertReplacerValue(t, repl, "extra.time.now.day", 2)
	assertReplacerValue(t, repl, "extra.time.now.day_padded", "02")
	assertReplacerValue(t, repl, "extra.time.now.hour", 3)
	assertReplacerValue(t, repl, "extra.time.now.hour_padded", "03")
	assertReplacerValue(t, repl, "extra.time.now.minute", 4)
	assertReplacerValue(t, repl, "extra.time.now.minute_padded", "04")
	assertReplacerValue(t, repl, "extra.time.now.second", 5)
	assertReplacerValue(t, repl, "extra.time.now.second_padded", "05")
	assertReplacerValue(t, repl, "extra.time.now.timezone_offset", "+0100")
	assertReplacerValue(t, repl, "extra.time.now.timezone_name", "CET")
	assertReplacerValue(t, repl, "extra.time.now.weekday_int", int(fixed.Weekday()))
	assertReplacerValue(t, repl, "extra.time.now.iso_week", isoWeek)
	assertReplacerValue(t, repl, "extra.time.now.iso_year", isoYear)
	assertReplacerValue(t, repl, "extra.time.now.custom", "2026/01/02 03:04")
}

func TestSetRandPlaceholders(t *testing.T) {
	e := ExtraPlaceholders{RandIntMin: 10, RandIntMax: 20}

	for range 200 {
		repl := caddy.NewReplacer()
		e.setRandPlaceholders(repl)

		floatVal, ok := repl.Get("extra.rand.float")
		if !ok {
			t.Fatal("missing extra.rand.float placeholder")
		}
		float64Val, ok := floatVal.(float64)
		if !ok {
			t.Fatalf("extra.rand.float has type %T, want float64", floatVal)
		}
		if float64Val < 0 || float64Val >= 1 {
			t.Fatalf("extra.rand.float = %v, want 0 <= value < 1", float64Val)
		}

		intVal, ok := repl.Get("extra.rand.int")
		if !ok {
			t.Fatal("missing extra.rand.int placeholder")
		}
		int64Val, ok := intVal.(int)
		if !ok {
			t.Fatalf("extra.rand.int has type %T, want int", intVal)
		}
		if int64Val < 10 || int64Val > 20 {
			t.Fatalf("extra.rand.int = %d, want 10 <= value <= 20", int64Val)
		}
	}
}

func TestSetLoadavgPlaceholders(t *testing.T) {
	oldLoadAvgFunc := loadAvgFunc
	loadAvgFunc = func() (*load.AvgStat, error) {
		return &load.AvgStat{
			Load1:  1.25,
			Load5:  2.5,
			Load15: 3.75,
		}, nil
	}
	t.Cleanup(func() {
		loadAvgFunc = oldLoadAvgFunc
	})

	repl := caddy.NewReplacer()
	ExtraPlaceholders{}.setLoadavgPlaceholders(repl)

	assertReplacerValue(t, repl, "extra.loadavg.1", 1.25)
	assertReplacerValue(t, repl, "extra.loadavg.5", 2.5)
	assertReplacerValue(t, repl, "extra.loadavg.15", 3.75)
}

func TestSetHostinfoPlaceholders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		oldHostUptimeFunc := hostUptimeFunc
		hostUptimeFunc = func() (uint64, error) {
			return 3661, nil
		}
		t.Cleanup(func() {
			hostUptimeFunc = oldHostUptimeFunc
		})

		repl := caddy.NewReplacer()
		ExtraPlaceholders{}.setHostinfoPlaceholders(repl)

		assertReplacerValue(t, repl, "extra.hostinfo.uptime", "1h1m1s")
	})

	t.Run("error", func(t *testing.T) {
		oldHostUptimeFunc := hostUptimeFunc
		hostUptimeFunc = func() (uint64, error) {
			return 0, errors.New("boom")
		}
		t.Cleanup(func() {
			hostUptimeFunc = oldHostUptimeFunc
		})

		repl := caddy.NewReplacer()
		ExtraPlaceholders{}.setHostinfoPlaceholders(repl)

		assertReplacerValue(t, repl, "extra.hostinfo.uptime", "error retrieving uptime")
	})
}

func TestServeHTTP(t *testing.T) {
	t.Run("sets placeholders and reuses the same instant", func(t *testing.T) {
		oldNowFunc := nowFunc
		oldHostUptimeFunc := hostUptimeFunc
		oldLoadAvgFunc := loadAvgFunc

		fixed := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.FixedZone("CET", 3600))
		loadAvgCalled := false

		nowFunc = func() time.Time { return fixed }
		hostUptimeFunc = func() (uint64, error) { return 90, nil }
		loadAvgFunc = func() (*load.AvgStat, error) {
			loadAvgCalled = true
			return &load.AvgStat{}, nil
		}

		t.Cleanup(func() {
			nowFunc = oldNowFunc
			hostUptimeFunc = oldHostUptimeFunc
			loadAvgFunc = oldLoadAvgFunc
		})

		repl := caddy.NewReplacer()
		req := httptest.NewRequest(http.MethodGet, "/path?q=a+b", nil)
		req.Host = "example.com"
		req = req.WithContext(context.WithValue(req.Context(), caddy.ReplacerCtxKey, repl))
		rr := httptest.NewRecorder()

		e := ExtraPlaceholders{
			RandIntMin:                 10,
			RandIntMax:                 20,
			TimeFormatCustom:           "2006-01-02T15:04:05Z07:00",
			DisableLoadavgPlaceholders: true,
		}

		nextCalled := false
		err := e.ServeHTTP(rr, req, caddyhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			nextCalled = true

			assertReplacerValue(t, repl, "extra.time.now.custom", "2026-01-02T03:04:05+01:00")
			assertReplacerValue(t, repl, "extra.time.now.utc.custom", "2026-01-02T02:04:05Z")
			assertReplacerValue(t, repl, "extra.hostinfo.uptime", "1m30s")
			assertReplacerValue(t, repl, "extra.newline", "\n")
			assertReplacerValue(t, repl, "extra.http.request.url.query_escaped", "http%3A%2F%2Fexample.com%2Fpath%3Fq%3Da%2Bb")

			if _, ok := repl.Get("extra.loadavg.1"); ok {
				t.Fatal("extra.loadavg.1 was set even though loadavg placeholders are disabled")
			}
			if mustGetString(t, repl, "extra.caddy.version.simple") == "" {
				t.Fatal("extra.caddy.version.simple is empty")
			}
			if mustGetString(t, repl, "extra.go.runtime.version") == "" {
				t.Fatal("extra.go.runtime.version is empty")
			}

			return nil
		}))
		if err != nil {
			t.Fatalf("ServeHTTP() returned error: %v", err)
		}
		if !nextCalled {
			t.Fatal("next handler was not called")
		}
		if loadAvgCalled {
			t.Fatal("loadAvgFunc was called even though loadavg placeholders are disabled")
		}
	})

	t.Run("returns internal server error without replacer", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rr := httptest.NewRecorder()

		nextCalled := false
		err := ExtraPlaceholders{}.ServeHTTP(rr, req, caddyhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			nextCalled = true
			return nil
		}))

		var handlerErr caddyhttp.HandlerError
		if !errors.As(err, &handlerErr) {
			t.Fatalf("ServeHTTP() error = %T, want caddyhttp.HandlerError", err)
		}
		if handlerErr.StatusCode != http.StatusInternalServerError {
			t.Fatalf("ServeHTTP() status code = %d, want %d", handlerErr.StatusCode, http.StatusInternalServerError)
		}
		if nextCalled {
			t.Fatal("next handler was called without replacer in context")
		}
	})
}

func assertReplacerValue(t *testing.T, repl *caddy.Replacer, key string, want any) {
	t.Helper()

	got, ok := repl.Get(key)
	if !ok {
		t.Fatalf("missing replacer key %q", key)
	}
	if got != want {
		t.Fatalf("%s = %#v, want %#v", key, got, want)
	}
}

func mustGetString(t *testing.T, repl *caddy.Replacer, key string) string {
	t.Helper()

	got, ok := repl.GetString(key)
	if !ok {
		t.Fatalf("missing replacer key %q", key)
	}
	return got
}
