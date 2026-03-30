# River Crossing

The river crossing problem has two kinds of participants:

- hackers
- serfs

A boat holds exactly four people. A crossing is safe if the boat contains:

- four hackers
- four serfs
- two hackers and two serfs

The unsafe combinations are one hacker with three serfs and one serf with
three hackers.

As each participant boards, it must invoke `board()`. All four participants in a
boatload must invoke `board()` before any participant in the next boatload
does. After all four have boarded, exactly one participant should invoke
`rowBoat()`.

In this package, a river crossing implementation must:

- allow only safe groups of four to board
- block incompatible or incomplete groups until a safe crew is available
- call `board()` for all four participants in a crew before any later crew boards
- call `rowBoat()` exactly once per completed crew
- allow waiting callers to stop waiting if their context is canceled
- be reusable across multiple crossings

## Layout

- `river_crossing.go`: challenge implementation
- `river_crossing_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable river crossing tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./river_crossing
```

Run the hint solutions:

```sh
go test ./river_crossing/solutions
```

Benchmark the hint solutions:

```sh
go test ./river_crossing/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./river_crossing -run '^$' -bench .
```
