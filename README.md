# tcrp
Timed Cache Refresh Proxy

This "proxy" will periodically perform a GET on a given URL and cache
the response body in memory which this proxy will then serve to any
GET requests it receives.

Internally, this uses sync/atomic to serve a consistent body to
clients.

Configuration occurs through environmental variables and this program
is otherwise stateless.

Variables:

`UPSTREAM_URL`: The upstream url to hit for content.
`REFRESH_PERIOD`: The frequency of updates (in seconds)
