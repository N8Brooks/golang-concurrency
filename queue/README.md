# Queue

The queue problem models leaders and followers waiting to dance in pairs.

When a leader arrives, it should pair with a waiting follower if one exists.
Otherwise it waits. Followers behave symmetrically.

In this package, a queue implementation must:

- let a matched leader and follower both proceed to dance
- block unmatched leaders and followers until a counterpart arrives
- allow a waiting caller to stop waiting if its context is canceled
- be reusable across multiple pairings

## Layout

- `queue.go`: challenge implementation stub
- `queue_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable queue tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./queue
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./queue/solutions
```

Benchmark the hint solutions:

```sh
go test ./queue/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./queue -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
