# Roller Coaster

The roller coaster problem models many passenger threads and a single car
thread. The car has a fixed capacity `C`, where `0 < C < n`, and it may depart
 only when exactly `C` passengers have boarded.

Each passenger does one ride:

1. wait for the car to load
2. board
3. wait for the car to unload
4. unboard

The car performs one trip:

1. load
2. wait for `C` passengers to board
3. run
4. unload
5. wait for the same `C` passengers to unboard

In this package, a roller coaster implementation must:

- prevent passengers from boarding before the car loads
- prevent the car from running until exactly `C` passengers have boarded
- prevent passengers from unboarding before the car unloads
- keep each trip isolated to one group of `C` passengers
- be reusable across many rides
- allow a waiting passenger or an empty car to stop waiting if its context is canceled

The provided hint solution uses a ride-by-ride handshake equivalent to the
textbook semaphore solution, with one waiting queue for passengers and a
barrier at boarding and unboarding.

## Layout

- `roller_coaster.go`: challenge implementation stub
- `roller_coaster_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable roller-coaster tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./roller_coaster
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./roller_coaster/solutions
```

Benchmark the hint solutions:

```sh
go test ./roller_coaster/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./roller_coaster -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
