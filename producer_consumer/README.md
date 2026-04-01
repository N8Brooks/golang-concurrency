# Producer Consumer

The producer-consumer problem has two kinds of participants:

1. producers, which wait for events and publish them
2. consumers, which receive published events and process them

The synchronization requirements are:

- a consumer must block while the buffer is empty
- `waitForEvent` and `process` should run outside the critical section
- a waiting consumer should return early if its context is canceled
- a canceled waiting consumer should not consume a future event
- the implementation should be reusable across many producer/consumer handoffs

Implementations in this repo may use either:

- a shared buffer protected by synchronization
- or a direct rendezvous/handoff between producers and consumers

The shared tests intentionally validate the synchronization contract rather than
requiring a specific internal queueing strategy.

The provided hint solutions illustrate three different approaches:

- `channel.go`: direct handoff via an unbuffered channel
- `sync_cond.go`: explicit shared buffer with `sync.Cond`
- `semaphore.go`: explicit shared buffer with semaphores

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
