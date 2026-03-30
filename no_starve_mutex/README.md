# No-Starve Mutex

The no-starve mutex problem is mutual exclusion with a bounded-waiting
guarantee.

Multiple threads need access to a critical section, but unlike a basic mutex, a
thread that starts waiting must not be overtaken indefinitely by later arrivals.

In this package, a no-starve mutex implementation must:

- allow only one caller into the critical section at a time
- prevent later arrivals from overtaking queued waiters indefinitely
- let a waiting caller stop waiting if its context is canceled
- be reusable across multiple acquire/release cycles

## Layout

- `no_starve_mutex.go`: challenge implementation stub
- `no_starve_mutex_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable no-starve-mutex tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./no_starve_mutex
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./no_starve_mutex/solutions
```

Benchmark the hint solutions:

```sh
go test ./no_starve_mutex/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./no_starve_mutex -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
