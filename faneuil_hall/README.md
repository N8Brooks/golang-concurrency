# Faneuil Hall

The Faneuil Hall problem has three kinds of participants:

- immigrants
- spectators
- one judge

Immigrants enter the building, check in, sit down, wait for confirmation,
swear, collect their certificates, and eventually leave. Spectators may enter,
watch, and leave. The judge enters, waits until all entered immigrants have
checked in, confirms the ceremony, and then leaves.

In this package, an implementation must:

- block immigrant and spectator entry while the judge is in the building
- block immigrant exit while the judge is in the building
- allow spectators to leave even while the judge is present
- prevent the judge from confirming until every entered immigrant has checked in
- prevent immigrants from getting certificates before the judge confirms
- support repeated ceremonies across multiple judge visits
- ensure that after a judge leaves, all sworn immigrants from that ceremony
  leave before the next judge may enter
- allow callers that are still waiting to enter to stop waiting if their
  context is canceled

## Layout

- `faneuil_hall.go`: challenge implementation stub
- `faneuil_hall_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./faneuil_hall
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./faneuil_hall/solutions
```

Benchmark the hint solutions:

```sh
go test ./faneuil_hall/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./faneuil_hall -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
