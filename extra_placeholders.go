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
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"go.uber.org/zap"
)

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
type ExtraPlaceholders struct {
	// RandIntMin defines the minimum value (inclusive) for the `{extra.rand.int}` placeholder.
	RandIntMin int `json:"rand_int_min,omitempty"`

	// RandIntMax defines the maximum value (inclusive) for the `{extra.rand.int}` placeholder.
	RandIntMax int `json:"rand_int_max,omitempty"`

	// TimeFormatCustom specifies a custom time format for the `{extra.time.now.custom}` placeholder.
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
		e.TimeFormatCustom = "2006-01-02 15:04:05" // Default format for custom time placeholder
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

	// Set placeholders for current time (month, day, hour, minute, second).
	now := time.Now() // System's local timezone
	repl.Set("extra.time.now.month", int(now.Month()))
	repl.Set("extra.time.now.month_padded", fmt.Sprintf("%02d", now.Month()))
	repl.Set("extra.time.now.day", now.Day())
	repl.Set("extra.time.now.day_padded", fmt.Sprintf("%02d", now.Day()))
	repl.Set("extra.time.now.hour", now.Hour())
	repl.Set("extra.time.now.hour_padded", fmt.Sprintf("%02d", now.Hour()))
	repl.Set("extra.time.now.minute", now.Minute())
	repl.Set("extra.time.now.minute_padded", fmt.Sprintf("%02d", now.Minute()))
	repl.Set("extra.time.now.second", now.Second())
	repl.Set("extra.time.now.second_padded", fmt.Sprintf("%02d", now.Second()))

	// Set placeholders for timezone offset and name.
	repl.Set("extra.time.now.timezone_offset", now.Format("-0700"))
	repl.Set("extra.time.now.timezone_name", now.Format("MST"))

	// Set placeholders for ISO week and ISO year.
	isoYear, isoWeek := now.ISOWeek()
	repl.Set("extra.time.now.iso_week", isoWeek)
	repl.Set("extra.time.now.iso_year", isoYear)

	// Set custom time format placeholder
	repl.Set("extra.time.now.custom", now.Format(e.TimeFormatCustom))

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
