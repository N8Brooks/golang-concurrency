# Finite Buffer Producer Consumer

This package models the finite-buffer variant of the producer-consumer problem.

Producers wait for events and add them to a shared buffer. Consumers remove
buffered events and process them.

The synchronization requirements are:

- access to the buffer must be exclusive while an item is added or removed
- a consumer must block while the buffer is empty
- a producer must block while the buffer is full
- `waitForEvent` and `process` should run outside the critical section
- waiting producers and consumers should return early if their context is canceled

In other words, this is the classic producer-consumer problem with both
`items` and `spaces` constraints.

## Layout

- `finite_buffer_producer_consumer.go`: challenge implementation
- `finite_buffer_producer_consumer_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable finite-buffer producer-consumer tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./finite_buffer_producer_consumer
```

Run the hint solutions:

```sh
go test ./finite_buffer_producer_consumer/solutions
```

Benchmark the hint solutions:

```sh
go test ./finite_buffer_producer_consumer/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./finite_buffer_producer_consumer -run '^$' -bench .
```
