package extraplaceholders

import (
	"runtime"
	"strconv"

	"github.com/caddyserver/caddy/v2"
)

// setGoPlaceholders sets placeholders for the Go runtime under extra.go.*
func (e ExtraPlaceholders) setGoPlaceholders(repl *caddy.Replacer) {
	repl.Set("extra.go.runtime.version", runtime.Version())
	repl.Set("extra.go.runtime.numgoroutines", strconv.Itoa(runtime.NumGoroutine()))
}
