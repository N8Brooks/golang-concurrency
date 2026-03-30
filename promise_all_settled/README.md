# Promise All Settled

This package models the Promise.allSettled style problem: run several
asynchronous functions in parallel and collect either the fulfilled value or the
rejection reason for each position.

## Layout

- `contract.go`: shared result type
- `promise_all_settled.go`: challenge stub
- `promise_all_settled_test.go`: shared tests for the challenge
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./promise_all_settled
```

Run the hint solutions:

```sh
go test ./promise_all_settled/solutions
```
