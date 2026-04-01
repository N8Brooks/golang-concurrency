# Promise All Settled

This package models the Promise.allSettled style problem: run several
asynchronous functions in parallel and return one result per position after all
of them settle.

Each position in the returned slice is a [`Result[T]`](./contract.go) with:

- `Val` set for a fulfilled promise
- `Err` set for a rejected promise

Unlike the promise-based helpers in this repo, `PromiseAllSettled` returns the
settled slice directly instead of wrapping it in another promise.

## Layout

- `contract.go`: shared `Result[T]` type
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

Benchmark the hint solutions:

```sh
go test ./promise_all_settled/solutions -run '^$' -bench .
```
