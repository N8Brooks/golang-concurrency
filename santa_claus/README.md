# Santa Claus

The Santa Claus problem has three kinds of participants:

- one sleeping Santa
- exactly nine reindeer
- any number of elves

Santa wakes up only when either:

- all nine reindeer have returned and the sleigh must be prepared, or
- exactly three elves are waiting for help

If both conditions hold, the reindeer take priority. Elves are helped in
batches of three, and no additional elves may enter the next batch until all
three elves in the current batch have finished `getHelp`.

In this package, an implementation must:

- call `prepareSleigh` once for each full batch of nine reindeer
- ensure all nine reindeer invoke `getHitched` only after `prepareSleigh`
- call `helpElves` once for each batch of three elves
- block additional elves until the current batch of three has finished
- keep Santa running in a loop so multiple elf batches can be served
- let waiting goroutines stop cleanly when their context is canceled

## Layout

- `santa_claus.go`: challenge implementation stub
- `santa_claus_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./santa_claus
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./santa_claus/solutions
```

Benchmark the hint solutions:

```sh
go test ./santa_claus/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./santa_claus -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
