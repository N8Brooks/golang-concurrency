# Cache With Time Limit

This package models a time-limited key/value cache. Entries expire after their
TTL, `Get` returns `-1` for missing keys, and `Count` reports the number of
currently live entries.

## Layout

- `cache_with_time_limit.go`: challenge stub
- `cache_with_time_limit_test.go`: shared tests for the challenge
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./cache_with_time_limit
```

Run the hint solutions:

```sh
go test ./cache_with_time_limit/solutions
```
