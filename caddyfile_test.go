package extraplaceholders

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func TestUnmarshalCaddyfile(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		d := caddyfile.NewTestDispenser(`extra_placeholders {
			rand_int 7 17
			time_format_custom 2006-01-02
			disable_loadavg_placeholders
		}`)

		var e ExtraPlaceholders
		if err := e.UnmarshalCaddyfile(d); err != nil {
			t.Fatalf("UnmarshalCaddyfile() returned error: %v", err)
		}

		if e.RandIntMin != 7 || e.RandIntMax != 17 {
			t.Fatalf("unexpected rand_int range: got %d..%d", e.RandIntMin, e.RandIntMax)
		}
		if e.TimeFormatCustom != "2006-01-02" {
			t.Fatalf("unexpected custom time format: got %q", e.TimeFormatCustom)
		}
		if !e.DisableLoadavgPlaceholders {
			t.Fatal("DisableLoadavgPlaceholders = false, want true")
		}
	})

	t.Run("invalid rand_int args", func(t *testing.T) {
		d := caddyfile.NewTestDispenser(`extra_placeholders {
			rand_int 7
		}`)

		var e ExtraPlaceholders
		if err := e.UnmarshalCaddyfile(d); err == nil {
			t.Fatal("UnmarshalCaddyfile() returned nil error, want error")
		}
	})

	t.Run("unknown subdirective", func(t *testing.T) {
		d := caddyfile.NewTestDispenser(`extra_placeholders {
			nope
		}`)

		var e ExtraPlaceholders
		if err := e.UnmarshalCaddyfile(d); err == nil {
			t.Fatal("UnmarshalCaddyfile() returned nil error, want error")
		}
	})
}
