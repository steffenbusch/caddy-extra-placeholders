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
	"math/rand"

	"github.com/caddyserver/caddy/v2"
)

// setRandPlaceholders sets placeholders for random float and integer values.
func (e ExtraPlaceholders) setRandPlaceholders(repl *caddy.Replacer) {
	repl.Set("extra.rand.float", rand.Float64())
	if e.RandIntMax > e.RandIntMin {
		repl.Set("extra.rand.int", rand.Intn(e.RandIntMax-e.RandIntMin+1)+e.RandIntMin)
	} else {
		repl.Set("extra.rand.int", rand.Intn(101)) // Default range 0-100 if not properly configured
	}
}
