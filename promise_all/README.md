# Promise All

This package models the Promise.all style problem: run a group of asynchronous
functions in parallel, preserve result order, and reject if any function fails.

## Layout

- `promise_all.go`: challenge stub
- `promise_all_test.go`: shared tests for the challenge
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./promise_all
```

Run the hint solutions:

```sh
go test ./promise_all/solutions
```
