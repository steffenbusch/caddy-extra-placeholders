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
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/shirou/gopsutil/v4/host"
)

// setHostinfoPlaceholders sets placeholders for system uptime in a human-readable format.
func (e ExtraPlaceholders) setHostinfoPlaceholders(repl *caddy.Replacer) {
	uptime, err := host.Uptime()
	if err == nil {
		uptimeDuration := time.Duration(uptime) * time.Second
		repl.Set("extra.hostinfo.uptime", uptimeDuration.String())
	} else {
		repl.Set("extra.hostinfo.uptime", "error retrieving uptime")
	}
}
