# Writers-Priority Readers-Writers

The writers-priority readers-writers problem models concurrent access to a
shared resource with an additional fairness requirement favoring writers.

Readers may access the resource at the same time, but writers require
exclusive access. Once a writer arrives, no new readers should be allowed to
enter until all queued writers have finished.

In this package, a writers-priority readers-writers implementation must:

- allow multiple readers into the critical section simultaneously
- prevent writers from entering while any reader is inside
- prevent any other reader or writer from entering while a writer is inside
- prevent newly arriving readers from entering while any writers are queued
- allow a waiting caller to stop waiting if its context is canceled
- remain reusable across multiple reads and writes

## Layout

- `writers_priority_readers_writers.go`: challenge implementation stub
- `writers_priority_readers_writers_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable writers-priority readers-writers tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./writers_priority_readers_writers
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./writers_priority_readers_writers/solutions
```

Benchmark the hint solutions:

```sh
go test ./writers_priority_readers_writers/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./writers_priority_readers_writers -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
