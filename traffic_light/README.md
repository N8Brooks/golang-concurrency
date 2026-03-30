# Traffic Light

The traffic light problem models a two-road intersection controlled by a
single light.

Road `1` starts green. When a car arrives:

- if its road is already green, it may cross immediately
- otherwise the light must be turned green for that road before the car crosses

In this package, a traffic-light implementation must:

- ensure every arriving car crosses exactly once
- ensure a car only crosses when its road is green
- avoid calling `turnGreen` when the correct road is already green
- remain reusable across many arriving cars

## Layout

- `traffic_light.go`: challenge implementation stub
- `traffic_light_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable traffic-light tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

The challenge implementation is opt-in and currently panics with
`panic("unimplemented")`.

Run the challenge tests:

```sh
go test -tags challenge ./traffic_light
```

That is expected to fail until the stub is implemented.

Run the hint solutions:

```sh
go test ./traffic_light/solutions
```

Benchmark the hint solutions:

```sh
go test ./traffic_light/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./traffic_light -run '^$' -bench .
```

Like the challenge tests, this is expected to fail until the stub is
implemented.
