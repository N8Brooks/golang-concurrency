# Bounded Blocking Queue

The bounded blocking queue problem models a queue with fixed capacity.

`Enqueue` blocks while the queue is full, `Dequeue` blocks while the queue is
empty, and `Size` reports the current number of elements.

In this package, an implementation must:

- block `Enqueue` when the queue is full
- block `Dequeue` when the queue is empty
- preserve all inserted values
- report the current size correctly
- be reusable across repeated enqueue/dequeue cycles

## Layout

- `bounded_blocking_queue.go`: challenge implementation stub
- `bounded_blocking_queue_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./bounded_blocking_queue
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./bounded_blocking_queue/solutions
```

Benchmark the hint solutions:

```sh
go test ./bounded_blocking_queue/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./bounded_blocking_queue -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
