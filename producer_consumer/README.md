# Producer Consumer

The producer-consumer problem has two kinds of participants:

1. producers, which wait for events and add them to a shared buffer
2. consumers, which remove buffered events and process them

The synchronization requirements are:

- access to the buffer must be exclusive while an item is added or removed
- a consumer must block while the buffer is empty
- `waitForEvent` and `process` should run outside the critical section
- a waiting consumer should return early if its context is canceled

In other words, producers and consumers coordinate through a shared buffer,
but they should not hold the buffer lock while waiting for new work or while
processing an item.

## Layout

- `producer_consumer.go`: challenge implementation
- `producer_consumer_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable producer-consumer tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./producer_consumer
```

Run the hint solutions:

```sh
go test ./producer_consumer/solutions
```

Benchmark the hint solutions:

```sh
go test ./producer_consumer/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./producer_consumer -run '^$' -bench .
```
