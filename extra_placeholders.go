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
	"net/http"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// defaultTimeFormatCustom is the fallback format used if no custom format is specified for the custom time placeholders.
const defaultTimeFormatCustom = "2006-01-02 15:04:05"

// ExtraPlaceholders provides additional placeholders that can be used within Caddy configurations:
//
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
//
// Current local time placeholders:
//
// Placeholder | Description
// ------------|-------------
// `{extra.time.now.month}` | Current month as an integer (e.g., 10 for October).
// `{extra.time.now.month_padded}` | Current month as a zero-padded string (e.g., "05" for May).
// `{extra.time.now.day}` | Current day of the month as an integer.
// `{extra.time.now.day_padded}` | Current day of the month as a zero-padded string.
// `{extra.time.now.hour}` | Current hour in 24-hour format as an integer.
// `{extra.time.now.hour_padded}` | Current hour in 24-hour format as a zero-padded string.
// `{extra.time.now.minute}` | Current minute as an integer.
// `{extra.time.now.minute_padded}` | Current minute as a zero-padded string.
// `{extra.time.now.second}` | Current second as an integer.
// `{extra.time.now.second_padded}` | Current second as a zero-padded string.
// `{extra.time.now.timezone_offset}` | Current timezone offset from UTC (e.g., +0200).
// `{extra.time.now.timezone_name}` | Current timezone abbreviation (e.g., CEST).
// `{extra.time.now.iso_week}` | Current ISO week number of the year.
// `{extra.time.now.iso_year}` | ISO year corresponding to the current ISO week.
// `{extra.time.now.custom}` | Current time in a custom format, configurable via the `time_format_custom` directive.
//
// UTC equivalents of the current time placeholders (with `.utc` added):
//
// Placeholder | Description
// ------------|-------------
// `{extra.time.now.utc.month}` | Current month in UTC as an integer (e.g., 10 for October).
// `{extra.time.now.utc.month_padded}` | Current month in UTC as a zero-padded string (e.g., "05" for May).
// `{extra.time.now.utc.day}` | Current day of the month in UTC as an integer.
// `{extra.time.now.utc.day_padded}` | Current day of the month in UTC as a zero-padded string.
// `{extra.time.now.utc.hour}` | Current hour in UTC in 24-hour format as an integer.
// `{extra.time.now.utc.hour_padded}` | Current hour in UTC in 24-hour format as a zero-padded string.
// `{extra.time.now.utc.minute}` | Current minute in UTC as an integer.
// `{extra.time.now.utc.minute_padded}` | Current minute in UTC as a zero-padded string.
// `{extra.time.now.utc.second}` | Current second in UTC as an integer.
// `{extra.time.now.utc.second_padded}` | Current second in UTC as a zero-padded string.
// `{extra.time.now.utc.timezone_offset}` | UTC timezone offset (always +0000).
// `{extra.time.now.utc.timezone_name}` | UTC timezone abbreviation (always UTC).
// `{extra.time.now.utc.iso_week}` | Current ISO week number of the year in UTC.
// `{extra.time.now.utc.iso_year}` | ISO year corresponding to the current ISO week in UTC.
// `{extra.time.now.utc.custom}` | Current UTC time in a custom format, configurable via the `time_format_custom` directive.
type ExtraPlaceholders struct {
	// RandIntMin defines the minimum value (inclusive) for the `{extra.rand.int}` placeholder.
	RandIntMin int `json:"rand_int_min,omitempty"`

	// RandIntMax defines the maximum value (inclusive) for the `{extra.rand.int}` placeholder.
	RandIntMax int `json:"rand_int_max,omitempty"`

	// TimeFormatCustom specifies a custom time format for the `{extra.time.now.custom}` and `{extra.time.now.utc.custom}` placeholder.
	// If left empty, a default format of "2006-01-02 15:04:05" is used.
	TimeFormatCustom string `json:"time_format_custom,omitempty"`

	// logger provides structured logging for the plugin's internal operations.
	logger *zap.Logger
}

// CaddyModule returns the module information required by Caddy to register the plugin.
func (ExtraPlaceholders) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.extra_placeholders",
		New: func() caddy.Module { return new(ExtraPlaceholders) },
	}
}

// Provision sets up the module. It is called once the module is instantiated.
func (e *ExtraPlaceholders) Provision(ctx caddy.Context) error {
	e.logger = ctx.Logger()

	// Set default values if not configured
	if e.RandIntMin == 0 && e.RandIntMax == 0 {
		e.RandIntMin = 0
		e.RandIntMax = 100
	}
	if e.TimeFormatCustom == "" {
		e.TimeFormatCustom = defaultTimeFormatCustom
	}

	// Log the chosen configuration values
	e.logger.Info("ExtraPlaceholders plugin configured",
		zap.Int("RandIntMin", e.RandIntMin),
		zap.Int("RandIntMax", e.RandIntMax),
		zap.String("TimeFormatCustom", e.TimeFormatCustom),
	)

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

	e.setCaddyPlaceholders(repl)
	e.setRandPlaceholders(repl)
	e.setLoadavgPlaceholders(repl)
	e.setHostinfoPlaceholders(repl)

	// Set time placeholders for server's local time
	e.setTimePlaceholders(repl, time.Now(), false)

	// Set time placeholders for UTC time
	e.setTimePlaceholders(repl, time.Now().UTC(), true)

	// Call the next handler in the chain.
	return next.ServeHTTP(w, r)
}

// Interface guards to ensure ExtraPlaceholders implements the necessary interfaces.
var (
	_ caddy.Module                = (*ExtraPlaceholders)(nil)
	_ caddy.Provisioner           = (*ExtraPlaceholders)(nil)
	_ caddy.Validator             = (*ExtraPlaceholders)(nil)
	_ caddyhttp.MiddlewareHandler = (*ExtraPlaceholders)(nil)
)
