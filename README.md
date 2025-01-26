# Extra Placeholders Caddy Plugin

The **Extra Placeholders** plugin for [Caddy](https://caddyserver.com) adds a variety of additional placeholders, enabling more dynamic and context-aware configurations. Built to enhance request matchers and response handling, this plugin allows Caddy to utilize real-time system metrics such as load averages, system uptime, and current Caddy version. It also provides versatile random number generation and a comprehensive range of time-based placeholders, including ISO week data, and customizable time formatting.

These placeholders unlock flexible, condition-driven logic within Caddy, letting configurations respond to server load, current time, or even random values. Whether used to influence routing decisions or enrich response data, this plugin empowers Caddy to meet specific operational needs by embedding live system insights directly into its configuration.

[![Go Report Card](https://goreportcard.com/badge/github.com/steffenbusch/caddy-extra-placeholders)](https://goreportcard.com/report/github.com/steffenbusch/caddy-extra-placeholders)

## Features

This plugin introduces new placeholders that can be used within Caddy configurations:

| Placeholder                          | Description                                           |
|--------------------------------------|-------------------------------------------------------|
| `{extra.caddy.version.simple}`       | Simple version information of the Caddy server (e.g., v2.8.4). |
| `{extra.caddy.version.full}`         | Full version information of the Caddy server (e.g., v2.8.4 h1:q3pe...k=). |
| `{extra.rand.float}`                 | Random float value between 0.0 and 1.0.               |
| `{extra.rand.int}`                   | Random integer value between the configured min and max (default is 0 to 100). |
| `{extra.loadavg.1}`                  | System load average over the last 1 minute.           |
| `{extra.loadavg.5}`                  | System load average over the last 5 minutes.          |
| `{extra.loadavg.15}`                 | System load average over the last 15 minutes.         |
| `{extra.hostinfo.uptime}`            | System uptime in a human-readable format.             |
| `{extra.newline}`                    | Newline character (\n).                               |

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
| `{extra.time.now.timezone_offset}`   | Current timezone offset from UTC (e.g., +0200).       |
| `{extra.time.now.timezone_name}`     | Current timezone abbreviation (e.g., CEST).           |
| `{extra.time.now.iso_week}`          | Current ISO week number of the year.                  |
| `{extra.time.now.iso_year}`          | ISO year corresponding to the current ISO week.       |
| `{extra.time.now.weekday_int}`       | Current day of the week as an integer (Sunday = 0, Monday = 1, ..., Saturday = 6). |
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
| `{extra.time.now.utc.timezone_name}` | UTC timezone abbreviation (always UTC).               |
| `{extra.time.now.utc.iso_week}`      | Current ISO week number of the year in UTC.           |
| `{extra.time.now.utc.iso_year}`      | ISO year corresponding to the current ISO week in UTC. |
| `{extra.time.now.utc.weekday_int}`   | Current day of the week in UTC as an integer (Sunday = 0, Monday = 1, ..., Saturday = 6). |
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

    respond "Caddy Version: {extra.caddy.version.full}, Uptime: {extra.hostinfo.uptime}, Random Int: {extra.rand.int}{extra.newline}"
}
```

This example demonstrates how to use the additional placeholders provided by this plugin to dynamically insert Caddy version information, system uptime, and a random integer between 10 and 50 into an HTTP response.

Using `{extra.newline}` at the end of the `respond` directive inserts a newline character. While it's possible to directly enter a newline in the Caddyfile, where the closing `"` would then be on a new line, inputting `"\n"` would not work as expected — it would be used literally in the response instead of as a newline. The `{extra.newline}` placeholder offers a clearer and more readable alternative for inserting actual newline characters.

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

#### Placeholder Support within `time_format_custom`

You can also specify placeholders within `time_format_custom`. For example, if you want the format to depend on an environment variable or request data, use `{env.*}` or `{http.request.*}` placeholders:

```caddyfile
extra_placeholders {
    time_format_custom {http.request.uri.query.format}
}
```

In this example, if a request contains a query parameter like `?format=02.01.2006`, the `{http.request.uri.query.format}` placeholder will dynamically resolve to `02.01.2006` for that specific request, which then sets the date format for `{extra.time.now.custom}`. This feature allows you to adjust time formatting based on user input, query parameters, or environment variables, creating a context-sensitive custom format.

> [!NOTE]
> When using placeholders in `time_format_custom`, ensure that the placeholder content aligns with [Go's time format syntax](https://pkg.go.dev/time#pkg-constants) to avoid formatting issues.

### Example: Conditional Redirect Based on Random Value

The following example demonstrates how you can use the [`map`](https://caddyserver.com/docs/caddyfile/directives/map) directive with the random integer placeholder to redirect users to different search engines based on the generated random number.
It also includes a response header to display the generated random integer, which can be useful for debugging purposes.

```caddyfile
{
    order extra_placeholders before header
}

:8080 {
    extra_placeholders {
        # Generate a random integer between 1 and 4, stored in {extra.rand.int},
        # matching the ranges defined in the "map" directive below
        rand_int 1 4
    }

    # Add the random integer as a response header
    header X-Random-Int {extra.rand.int}

    # Map the random integer value to the redirection target
    map {extra.rand.int} {redir_target} {
        1 https://www.google.com
        2 https://www.bing.com
        3 https://www.yahoo.com
        4 https://www.duckduckgo.com
    }

    # Redirect to the mapped target
    redir {redir_target}
}
```

> [!NOTE]
> In this example, the global `order` directive is used to overwrite the default directive order of the `extra_placeholders` directive, which is ordered before `redir` in the Caddyfile.
> The `header` directive is processed even earlier than `redir` in the directive sorting order. To ensure that the `header` directive has the necessary placeholder values for processing, the `extra_placeholders` directive must be evaluated before `header`

### Example: Time-Based Greeting

The following example demonstrates how you can use conditional expressions with the `extra.time.now.hour` placeholder to greet users with a time-appropriate message:

```caddyfile
:8080 {
    extra_placeholders

    @morning `{extra.time.now.hour} >= 6 && {extra.time.now.hour} <12`
    @day `{extra.time.now.hour} >= 12 && {extra.time.now.hour} <18`
    @evening `{extra.time.now.hour} >= 18 && {extra.time.now.hour} <22`
    @night `{extra.time.now.hour} >= 22 || {extra.time.now.hour} <6`

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

> [!NOTE]
> If the first argument of a named matcher starts with a quoted token, it is automatically treated as an expression, making the `expression` keyword unnecessary.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Caddy](https://caddyserver.com) for providing a powerful and extensible web server.
- [gopsutil](https://github.com/shirou/gopsutil) for system metrics, such as load averages and uptime, used under the BSD 3-Clause License.
