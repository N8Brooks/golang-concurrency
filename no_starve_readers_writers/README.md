# No-Starve Readers-Writers

The no-starve readers-writers problem models concurrent access to a shared
resource with an additional fairness requirement for writers.

Readers may access the resource at the same time, but writers require
exclusive access. Once a writer is queued, newly arriving readers should not be
allowed to pass ahead of it indefinitely.

In this package, a no-starve readers-writers implementation must:

- allow multiple readers into the critical section simultaneously
- prevent writers from entering while any reader is inside
- prevent any other reader or writer from entering while a writer is inside
- prevent newly arriving readers from bypassing a queued writer
- allow a waiting caller to stop waiting if its context is canceled
- remain reusable across multiple reads and writes

## Layout

- `no_starve_readers_writers.go`: challenge implementation stub
- `no_starve_readers_writers_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable no-starve readers-writers tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./no_starve_readers_writers
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./no_starve_readers_writers/solutions
```

Benchmark the hint solutions:

```sh
go test ./no_starve_readers_writers/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./no_starve_readers_writers -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
