package extraplaceholders

import (
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
)

// setTimePlaceholders sets placeholders for date, time, and custom format,
// using the provided time.Time. If isUTC is true, ".utc" is added in the placeholder path.
func (e ExtraPlaceholders) setTimePlaceholders(repl *caddy.Replacer, t time.Time, isUTC bool) {
	// Determine the base path, with or without ".utc"
	base := "extra.time.now"
	if isUTC {
		base += ".utc"
	}

	// Set date and time components with the specified base path
	repl.Set(fmt.Sprintf("%s.month", base), int(t.Month()))
	repl.Set(fmt.Sprintf("%s.month_padded", base), fmt.Sprintf("%02d", t.Month()))
	repl.Set(fmt.Sprintf("%s.day", base), t.Day())
	repl.Set(fmt.Sprintf("%s.day_padded", base), fmt.Sprintf("%02d", t.Day()))
	repl.Set(fmt.Sprintf("%s.hour", base), t.Hour())
	repl.Set(fmt.Sprintf("%s.hour_padded", base), fmt.Sprintf("%02d", t.Hour()))
	repl.Set(fmt.Sprintf("%s.minute", base), t.Minute())
	repl.Set(fmt.Sprintf("%s.minute_padded", base), fmt.Sprintf("%02d", t.Minute()))
	repl.Set(fmt.Sprintf("%s.second", base), t.Second())
	repl.Set(fmt.Sprintf("%s.second_padded", base), fmt.Sprintf("%02d", t.Second()))

	// Set timezone offset and name
	repl.Set(fmt.Sprintf("%s.timezone_offset", base), t.Format("-0700"))
	repl.Set(fmt.Sprintf("%s.timezone_name", base), t.Format("MST"))

	// Set ISO week and year components
	isoYear, isoWeek := t.ISOWeek()
	repl.Set(fmt.Sprintf("%s.iso_week", base), isoWeek)
	repl.Set(fmt.Sprintf("%s.iso_year", base), isoYear)

	// Set custom time format placeholder
	repl.Set(fmt.Sprintf("%s.custom", base), t.Format(e.TimeFormatCustom))
}
