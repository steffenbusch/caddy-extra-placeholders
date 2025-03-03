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
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	// Register the module with Caddy and specify where in the directive order it should be applied.
	caddy.RegisterModule(ExtraPlaceholders{})
	httpcaddyfile.RegisterHandlerDirective("extra_placeholders", parseCaddyfile)
	httpcaddyfile.RegisterDirectiveOrder("extra_placeholders", "before", "redir")
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
		case "time_format_custom":
			if d.NextArg() {
				e.TimeFormatCustom = d.Val()
			} else {
				return d.ArgErr()
			}
		case "disable_loadavg_placeholders":
			e.DisableLoadavgPlaceholders = true
		default:
			// Handle unknown subdirective with an error message
			return d.Errf("unknown subdirective: %s", d.Val())
		}
	}
	return nil
}
