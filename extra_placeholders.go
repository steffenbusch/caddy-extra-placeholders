// Copyright 2024 Steffen Busch

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package extraplaceholders

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
)

// ExtraPlaceholders represents the structure for the plugin.
type ExtraPlaceholders struct {
	RandIntMin int `json:"rand_int_min,omitempty"`
	RandIntMax int `json:"rand_int_max,omitempty"`
}

func init() {
	// Register the module with Caddy and specify where in the directive order it should be applied.
	caddy.RegisterModule(ExtraPlaceholders{})
	httpcaddyfile.RegisterHandlerDirective("extra_placeholders", parseCaddyfile)
	httpcaddyfile.RegisterDirectiveOrder("extra_placeholders", "before", "redir")
}

// Placeholder | Description
// ------------|-------------
// `{extra.caddy.version.simple}` | Simple version information of the Caddy server.
// `{extra.caddy.version.full}` | Full version information of the Caddy server.
// `{extra.rand.float}` | Random float value between 0.0 and 1.0.
// `{extra.rand.int}` | Random integer value between the configured min and max (default is 0 to 100).
// `{extra.loadavg.1}` | System load average over the last 1 minute.
// `{extra.loadavg.5}` | System load average over the last 5 minutes.
// `{extra.loadavg.15}` | System load average over the last 15 minutes.
// `{extra.hostinfo.uptime}` | System uptime in a human-readable format.

// CaddyModule returns the module information required by Caddy to register the plugin.
func (ExtraPlaceholders) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.extra_placeholders",
		New: func() caddy.Module { return new(ExtraPlaceholders) },
	}
}

// Provision sets up the module. It is called once the module is instantiated.
func (e *ExtraPlaceholders) Provision(ctx caddy.Context) error {
	// Set default values if not configured
	if e.RandIntMin == 0 && e.RandIntMax == 0 {
		e.RandIntMin = 0
		e.RandIntMax = 100
	}
	return nil
}

// Validate ensures the configuration is correct.
func (e *ExtraPlaceholders) Validate() error {
	if e.RandIntMax <= e.RandIntMin {
		return fmt.Errorf("invalid configuration: RandIntMax (%d) must be greater than RandIntMin (%d)", e.RandIntMax, e.RandIntMin)
	}
	return nil
}

// ServeHTTP adds new placeholders and passes the request to the next handler in the chain.
func (e ExtraPlaceholders) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Retrieve the replacer from the request context.
	repl, ok := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	if !ok {
		return caddyhttp.Error(http.StatusInternalServerError, nil)
	}

	// Set Caddy version placeholders.
	simpleVersion, fullVersion := caddy.Version()
	repl.Set("extra.caddy.version.simple", simpleVersion)
	repl.Set("extra.caddy.version.full", fullVersion)

	// Set placeholders for random float and integer values.
	repl.Set("extra.rand.float", rand.Float64())
	if e.RandIntMax > e.RandIntMin {
		repl.Set("extra.rand.int", rand.Intn(e.RandIntMax-e.RandIntMin+1)+e.RandIntMin)
	} else {
		repl.Set("extra.rand.int", rand.Intn(101)) // Default range 0-100 if not properly configured
	}

	// Set placeholders for system load averages (1, 5, and 15 minutes).
	loadAvg, err := load.Avg()
	if err == nil {
		repl.Set("extra.loadavg.1", loadAvg.Load1)
		repl.Set("extra.loadavg.5", loadAvg.Load5)
		repl.Set("extra.loadavg.15", loadAvg.Load15)
	}

	// Set placeholder for system uptime.
	uptime, err := host.Uptime()
	if err == nil {
		uptimeDuration := time.Duration(uptime) * time.Second
		repl.Set("extra.hostinfo.uptime", uptimeDuration.String())
	} else {
		repl.Set("extra.hostinfo.uptime", "error retrieving uptime")
	}

	// Call the next handler in the chain.
	return next.ServeHTTP(w, r)
}

// parseCaddyfile parses tokens from the Caddyfile into a new ExtraPlaceholders instance.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m ExtraPlaceholders
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// UnmarshalCaddyfile processes the configuration from the Caddyfile.
func (e *ExtraPlaceholders) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	// Consume the directive name.
	d.Next()

	for d.NextBlock(0) {
		switch d.Val() {
		case "rand_int":
			args := d.RemainingArgs()
			if len(args) != 2 {
				return d.ArgErr()
			}
			min, err1 := strconv.Atoi(args[0])
			max, err2 := strconv.Atoi(args[1])
			if err1 != nil || err2 != nil {
				return d.ArgErr()
			}
			e.RandIntMin = min
			e.RandIntMax = max
		}
	}
	return nil
}

// Interface guards to ensure ExtraPlaceholders implements the necessary interfaces.
var (
	_ caddy.Module                = (*ExtraPlaceholders)(nil)
	_ caddy.Provisioner           = (*ExtraPlaceholders)(nil)
	_ caddy.Validator             = (*ExtraPlaceholders)(nil)
	_ caddyhttp.MiddlewareHandler = (*ExtraPlaceholders)(nil)
)
