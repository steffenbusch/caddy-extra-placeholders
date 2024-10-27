# Extra Placeholders Caddy Plugin

The **Extra Placeholders** plugin for [Caddy](https://caddyserver.com) adds a variety of additional placeholders, enabling more dynamic and context-aware configurations. Built to enhance request matchers and response handling, this plugin allows Caddy to utilize real-time system metrics such as load averages, system uptime, and current Caddy version. It also provides versatile random number generation and a comprehensive range of time-based placeholders, including ISO week data, and customizable time formatting.

These placeholders unlock flexible, condition-driven logic within Caddy, letting configurations respond to server load, current time, or even random values. Whether used to influence routing decisions or enrich response data, this plugin empowers Caddy to meet specific operational needs by embedding live system insights directly into its configuration.

## Features

This plugin introduces new placeholders that can be used within Caddy configurations:

| Placeholder                          | Description                                           |
|--------------------------------------|-------------------------------------------------------|
| `{extra.caddy.version.simple}`       | Simple version information of the Caddy server.       |
| `{extra.caddy.version.full}`         | Full version information of the Caddy server.         |
| `{extra.rand.float}`                 | Random float value between 0.0 and 1.0.               |
| `{extra.rand.int}`                   | Random integer value between the configured min and max (default is 0 to 100). |
| `{extra.loadavg.1}`                  | System load average over the last 1 minute.           |
| `{extra.loadavg.5}`                  | System load average over the last 5 minutes.          |
| `{extra.loadavg.15}`                 | System load average over the last 15 minutes.         |
| `{extra.hostinfo.uptime}`            | System uptime in a human-readable format.             |

### Current Server Local Time Placeholders

These placeholders reflect the **server's local timezone**:

| Placeholder                          | Description                                           |
|--------------------------------------|-------------------------------------------------------|
| `{extra.time.now.month}`             | Current month as an integer (e.g., 5 for May).        |
| `{extra.time.now.month_padded}`      | Current month as a zero-padded string (e.g., "05" for May). |
| `{extra.time.now.day}`               | Current day of the month as an integer.               |
| `{extra.time.now.day_padded}`        | Current day of the month as a zero-padded string.     |
| `{extra.time.now.hour}`              | Current hour in 24-hour format as an integer.         |
| `{extra.time.now.hour_padded}`       | Current hour in 24-hour format as a zero-padded string. |
| `{extra.time.now.minute}`            | Current minute as an integer.                         |
| `{extra.time.now.minute_padded}`     | Current minute as a zero-padded string.               |
| `{extra.time.now.second}`            | Current second as an integer.                         |
| `{extra.time.now.second_padded}`     | Current second as a zero-padded string.               |
| `{extra.time.now.timezone_offset}`   | The current timezone offset from UTC (e.g., +0200).   |
| `{extra.time.now.timezone_name}`     | The current timezone abbreviation (e.g., CEST).       |
| `{extra.time.now.iso_week}`          | The current ISO week number of the year.              |
| `{extra.time.now.iso_year}`          | The ISO year corresponding to the current ISO week.   |
| `{extra.time.now.custom}`            | Current time in a custom format, configurable via the `time_format_custom` directive. |

### Current UTC Time Placeholders

These placeholders reflect **UTC time**:

| Placeholder                          | Description                                           |
|--------------------------------------|-------------------------------------------------------|
| `{extra.time.now.utc.month}`         | Current month in UTC as an integer (e.g., 5 for May). |
| `{extra.time.now.utc.month_padded}`  | Current month in UTC as a zero-padded string (e.g., "05" for May). |
| `{extra.time.now.utc.day}`           | Current day of the month in UTC as an integer.        |
| `{extra.time.now.utc.day_padded}`    | Current day of the month in UTC as a zero-padded string. |
| `{extra.time.now.utc.hour}`          | Current hour in UTC in 24-hour format as an integer.  |
| `{extra.time.now.utc.hour_padded}`   | Current hour in UTC in 24-hour format as a zero-padded string. |
| `{extra.time.now.utc.minute}`        | Current minute in UTC as an integer.                  |
| `{extra.time.now.utc.minute_padded}` | Current minute in UTC as a zero-padded string.        |
| `{extra.time.now.utc.second}`        | Current second in UTC as an integer.                  |
| `{extra.time.now.utc.second_padded}` | Current second in UTC as a zero-padded string.        |
| `{extra.time.now.utc.timezone_offset}` | UTC timezone offset (always +0000).                 |
| `{extra.time.now.utc.timezone_name}` | UTC timezone abbreviation (always UTC).              |
| `{extra.time.now.utc.iso_week}`      | Current ISO week number of the year in UTC.           |
| `{extra.time.now.utc.iso_year}`      | ISO year corresponding to the current ISO week in UTC. |
| `{extra.time.now.utc.custom}`        | Current UTC time in a custom format, configurable via the `time_format_custom` directive. |

> [!NOTE]
> All `extra.time.now.*` placeholders refer to the system's local timezone, while `extra.time.now.utc.*` placeholders represent the same values in UTC.

## Building

To build Caddy with this module, use [xcaddy](https://github.com/caddyserver/xcaddy):

```bash
$ xcaddy build --with github.com/steffenbusch/caddy-extra-placeholders
```

## Caddyfile config

By default, the `extra_placeholders` directive is ordered before `redir` in the Caddyfile. This simplifies configuration and removes the need for manual ordering in most cases.
However, if this default order does not suit your needs, you can still change the order using the [`order`](https://caddyserver.com/docs/caddyfile/directives#directive-order) global directive.

To use the extra placeholders, you can add the following directive to your Caddyfile:

```caddyfile
:8080 {
    extra_placeholders {
        rand_int 10 50
    }

    respond "Caddy Version: {extra.caddy.version.full}, Uptime: {extra.hostinfo.uptime}, Random Int: {extra.rand.int}"
}
```

This example demonstrates how to use the additional placeholders provided by this plugin to dynamically insert Caddy version information, system uptime, and a random integer between 10 and 50 into an HTTP response.

### Random Integer Configuration

To configure the range for the `{extra.rand.int}` placeholder, use the `rand_int` subdirective inside the `extra_placeholders` directive. The format is:

```caddyfile
extra_placeholders {
    rand_int <min> <max>
}
```

- `<min>`: The minimum value for the random integer (inclusive).
- `<max>`: The maximum value for the random integer (inclusive).

If `rand_int` is not specified, the default values are:

- `RandIntMin = 0`
- `RandIntMax = 100`

This means that `{extra.rand.int}` will default to generating a random integer between 0 and 100 if not explicitly configured.

### Custom Time Format

The `{extra.time.now.custom}` and `{extra.time.now.utc.custom}` placeholders can be configured using the `time_format_custom` subdirective inside the `extra_placeholders` directive.
This allows you to specify a custom date and time format using [Go's time format syntax](https://pkg.go.dev/time#pkg-constants).

```caddyfile
extra_placeholders {
    # For DD.MM.YYYY HH:MM:SS
    time_format_custom "02.01.2006 15:04:05"
}
```

If `time_format_custom` is not specified, it defaults to `"2006-01-02 15:04:05"`. This format will be applied to both `{extra.time.now.custom}` (server’s local timezone) and `{extra.time.now.utc.custom}` (UTC time) placeholders.

### Example: Conditional Redirect Based on Random Value

The following example demonstrates how you can use conditional expressions with the random integer placeholder to redirect users to different search engines based on the generated random number:

```caddyfile
:8080 {
    extra_placeholders {
        rand_int 1 100
    }

    @redirectToGoogle expression `{extra.rand.int} <= 25`
    redir @redirectToGoogle https://www.google.com

    @redirectToBing expression `{extra.rand.int} > 25 && {extra.rand.int} <= 50`
    redir @redirectToBing https://www.bing.com

    @redirectToYahoo expression `{extra.rand.int} > 50 && {extra.rand.int} <= 75`
    redir @redirectToYahoo https://www.yahoo.com

    @redirectToDuckDuckGo expression `{extra.rand.int} > 75`
    redir @redirectToDuckDuckGo https://www.duckduckgo.com
}
```

In this example:

- If `{extra.rand.int}` is between 1 and 25, the request is redirected to **Google** ([https://www.google.com](https://www.google.com)).
- If `{extra.rand.int}` is between 26 and 50, the request is redirected to **Bing** ([https://www.bing.com](https://www.bing.com)).
- If `{extra.rand.int}` is between 51 and 75, the request is redirected to **Yahoo** ([https://www.yahoo.com](https://www.yahoo.com)).
- If `{extra.rand.int}` is greater than 75, the request is redirected to **DuckDuckGo** ([https://www.duckduckgo.com](https://www.duckduckgo.com)).

This example demonstrates how to use the random integer placeholder in combination with conditional expressions to create dynamic redirection rules.

### Example: Time-Based Greeting

The following example demonstrates how you can use conditional expressions with the `extra.time.now.hour` placeholder to greet users with a time-appropriate message:

```caddyfile
:8080 {
    extra_placeholders

    @morning {
        expression `{extra.time.now.hour} >= 6 && {extra.time.now.hour} < 12`
    }
    @day {
        expression `{extra.time.now.hour} >= 12 && {extra.time.now.hour} < 18`
    }
    @evening {
        expression `{extra.time.now.hour} >= 18 && {extra.time.now.hour} < 22`
    }
    @night {
        expression `{extra.time.now.hour} >= 22 || {extra.time.now.hour} < 6`
    }

    handle @morning {
        respond "Good morning! And in case I don't see ya, good afternoon, good evening, and good night!"
    }

    handle @day {
        respond "Good day!"
    }

    handle @evening {
        respond "Good evening!"
    }

    handle @night {
        respond "Good night!"
    }
}
```

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Caddy](https://caddyserver.com) for providing a powerful and extensible web server.
- [gopsutil](https://github.com/shirou/gopsutil) for system metrics, such as load averages and uptime, used under the BSD 3-Clause License.
