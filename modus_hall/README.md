# Modus Hall

The Modus Hall problem models two factions trying to cross a narrow field:

- heathens
- prudes

Members of the same faction may cross concurrently, but members of opposite
factions may not cross at the same time.

Control of the field is determined by majority rule. While one faction is
crossing, the other faction may accumulate in queue. Once the queued faction
strictly outnumbers the active faction, new entrants from the active faction
must be blocked. Control then flips when the field becomes empty.

In this package, an implementation must:

- allow multiple heathens to cross together
- allow multiple prudes to cross together
- never allow heathens and prudes to cross simultaneously
- keep admitting the active faction until the queued opposition gains a strict majority
- block new entrants from the active faction once a transition begins
- switch control to the queued majority when the field clears
- allow waiting callers to stop waiting if their context is canceled
- be reusable across many control handoffs

## Layout

- `modus_hall.go`: challenge implementation stub
- `modus_hall_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./modus_hall
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./modus_hall/solutions
```

Benchmark the hint solutions:

```sh
go test ./modus_hall/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./modus_hall -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
