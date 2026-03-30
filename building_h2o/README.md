# Building H2O

The building H2O problem has two kinds of threads: hydrogen and oxygen.

Water molecules require one oxygen and two hydrogens. Threads should wait
until a complete molecule is ready, and all three threads in one molecule must
invoke `bond` before any thread from the next molecule does.

In this package, an H2O implementation must:

- allow exactly two hydrogens and one oxygen to bond per molecule
- block incomplete groups until a full molecule can be formed
- ensure bond invocations occur in complete molecules
- remain reusable across multiple molecules

## Layout

- `building_h2o.go`: challenge implementation stub
- `building_h2o_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable H2O tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./building_h2o
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./building_h2o/solutions
```

Benchmark the hint solutions:

```sh
go test ./building_h2o/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./building_h2o -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
