# Readers-Writers

The readers-writers problem models concurrent access to a shared resource.

Readers may access the resource at the same time, but writers require
exclusive access.

In this package, a readers-writers implementation must:

- allow multiple readers into the critical section simultaneously
- prevent writers from entering while any reader is inside
- prevent any other reader or writer from entering while a writer is inside
- allow a waiting caller to stop waiting if its context is canceled
- remain reusable across multiple reads and writes

## Layout

- `readers_writers.go`: challenge implementation stub
- `readers_writers_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable readers-writers tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./readers_writers
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./readers_writers/solutions
```

Benchmark the hint solutions:

```sh
go test ./readers_writers/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./readers_writers -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
