# Multi-Car Roller Coaster

The multi-car roller coaster problem extends the roller coaster problem to `m`
cars, each with capacity `C`.

Each car has a stable identifier in `[0, m)`. Cars take turns entering the
loading area and, because they cannot pass each other on the track, they must
also unload in that same cyclic order.

Each passenger does one ride:

1. wait for some car to load
2. board
3. wait for that car to unload
4. unboard

Each car performs one trip:

1. wait for its turn to load
2. load
3. wait for `C` passengers to board
4. release the next car to begin loading
5. run
6. wait for its turn to unload
7. unload
8. wait for its `C` passengers to unboard
9. release the next car to unload

In this package, a multi-car roller coaster implementation must:

- allow only one car to board passengers at a time
- allow multiple cars to be on the track concurrently
- unload cars in the same order they boarded
- ensure one carload finishes unboarding before the next carload begins to unboard
- be reusable across many trips
- allow a waiting passenger or not-yet-filled car to stop waiting if its context is canceled

The provided hint solution mirrors the textbook structure with one loading-area
turn and one unloading-area turn, both cycled by car identifier.

## Layout

- `multi_car_roller_coaster.go`: challenge implementation stub
- `multi_car_roller_coaster_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable multi-car roller-coaster tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./multi_car_roller_coaster
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./multi_car_roller_coaster/solutions
```

Benchmark the hint solutions:

```sh
go test ./multi_car_roller_coaster/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./multi_car_roller_coaster -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
