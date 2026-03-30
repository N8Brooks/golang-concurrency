# Web Crawler

The web crawler problem starts from a URL and an `HTMLParser`.

The crawler should:

- visit the start URL
- follow only URLs on the same host
- avoid revisiting URLs it has already seen
- return every reachable same-host URL

This version is intentionally concurrency-oriented: implementations may crawl
multiple pages in parallel.

## Layout

- `contract.go`: shared crawler and parser interfaces
- `web_crawler.go`: challenge implementation stub
- `web_crawler_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable crawler tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./web_crawler
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./web_crawler/solutions
```

Benchmark the hint solutions:

```sh
go test ./web_crawler/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./web_crawler -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
