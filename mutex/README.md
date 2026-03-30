# Mutex

The mutex problem is the basic mutual exclusion problem.

Multiple threads need access to a critical section, but at most one thread may
be in that critical section at a time.

In this package, a mutex implementation must:

- allow one caller to acquire the mutex
- block other callers until it is released
- let a waiting caller stop waiting if its context is canceled
- be reusable across multiple acquire/release cycles

## Layout

- `mutex.go`: challenge implementation stub
- `mutex_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable mutex tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./mutex
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./mutex/solutions
```

Benchmark the hint solutions:

```sh
go test ./mutex/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./mutex -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
