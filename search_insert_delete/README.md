# Search Insert Delete

The search-insert-delete problem models three different categories of access to
the same linked list:

- searchers, which only examine the list
- inserters, which append to the end of the list
- deleters, which remove items from anywhere in the list

The synchronization requirements are:

- searchers may run concurrently with other searchers
- inserters must be mutually exclusive with other inserters
- inserters may run concurrently with any number of searchers
- deleters must be mutually exclusive with searchers, inserters, and other deleters
- waiting callers should be able to stop waiting if their context is canceled

In other words, this is a three-way categorical exclusion problem: searchers
share with searchers and inserters, inserters share with searchers but not with
other inserters, and deleters share with nobody.

## Layout

- `search_insert_delete.go`: challenge implementation
- `search_insert_delete_test.go`: shared tests run against the challenge implementation
- `testsuite/`: reusable search-insert-delete tests and benchmarks
- `solutions/`: hint implementations that also use the shared testsuite

## Testing

Run the challenge tests:

```sh
go test -tags challenge ./search_insert_delete
```

Run the hint solutions:

```sh
go test ./search_insert_delete/solutions
```

Benchmark the hint solutions:

```sh
go test ./search_insert_delete/solutions -run '^$' -bench .
```

Benchmark the challenge implementation:

```sh
go test -tags challenge ./search_insert_delete -run '^$' -bench .
```
