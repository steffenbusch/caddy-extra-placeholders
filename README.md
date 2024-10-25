# Extra Placeholders Caddy Plugin

This repository contains the **Extra Placeholders** plugin for the [Caddy](https://caddyserver.com) server. This plugin provides additional placeholders that can be used for enhanced server-side information, such as system load averages, random numbers, Caddy version details, and system uptime.

## Features

This plugin introduces new placeholders that can be used within Caddy configurations:

| Placeholder                    | Description                                       |
|--------------------------------|---------------------------------------------------|
| `{extra.caddy.version.simple}` | Simple version information of the Caddy server.   |
| `{extra.caddy.version.full}`   | Full version information of the Caddy server.     |
| `{extra.rand.float}`           | Random float value between 0.0 and 1.0.           |
| `{extra.rand.int.0-100}`       | Random integer value between 0 and 100.           |
| `{extra.load1}`                | System load average over the last 1 minute.       |
| `{extra.load5}`                | System load average over the last 5 minutes.      |
| `{extra.load15}`               | System load average over the last 15 minutes.     |
| `{extra.uptime}`               | System uptime in a human-readable format.         |

These placeholders can be used in Caddyfiles to provide dynamic content and system information in responses.

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
    extra_placeholders

    respond "Caddy Version: {extra.caddy.version.full}, Uptime: {extra.uptime}"
}
```

This example demonstrates how to use the additional placeholders provided by this plugin to dynamically insert Caddy version information and system uptime into an HTTP response.

### Example: Conditional Redirect Based on Random Value

The following example demonstrates how you can use conditional expressions with the random integer placeholder to redirect users to different search engines based on the generated random number:

```caddyfile
:8080 {
    extra_placeholders

    @redirectToGoogle expression `{extra.rand.int.0-100} <= 25`
    redir @redirectToGoogle https://www.google.com

    @redirectToBing expression `{extra.rand.int.0-100} > 25 && {extra.rand.int.0-100} <= 50`
    redir @redirectToBing https://www.bing.com

    @redirectToYahoo expression `{extra.rand.int.0-100} > 50 && {extra.rand.int.0-100} <= 75`
    redir @redirectToYahoo https://www.yahoo.com

    @redirectToDuckDuckGo expression `{extra.rand.int.0-100} > 75`
    redir @redirectToDuckDuckGo https://www.duckduckgo.com
}
```

In this example:

- If `{extra.rand.int.0-100}` is between 0 and 25, the request is redirected to **Google** ([https://www.google.com](https://www.google.com)).
- If `{extra.rand.int.0-100}` is between 26 and 50, the request is redirected to **Bing** ([https://www.bing.com](https://www.bing.com)).
- If `{extra.rand.int.0-100}` is between 51 and 75, the request is redirected to **Yahoo** ([https://www.yahoo.com](https://www.yahoo.com)).
- If `{extra.rand.int.0-100}` is greater than 75, the request is redirected to **DuckDuckGo** ([https://www.duckduckgo.com](https://www.duckduckgo.com)).

This example demonstrates how to use the random integer placeholder in combination with conditional expressions to create dynamic redirection rules.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Caddy](https://caddyserver.com) for providing a powerful and extensible web server.
- [gopsutil](https://github.com/shirou/gopsutil) for system metrics, such as load averages and uptime, used under the BSD 3-Clause License.
