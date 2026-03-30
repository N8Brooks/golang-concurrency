# Promise Time Limit

This package wraps an asynchronous integer function with a deadline. If the
underlying work does not finish in time, the returned promise rejects with
`time limit exceeded`.

## Layout

- `promise_time_limit.go`: challenge stub
- `promise_time_limit_test.go`: shared tests for the challenge
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./promise_time_limit
```

Run the hint solutions:

```sh
go test ./promise_time_limit/solutions
```
