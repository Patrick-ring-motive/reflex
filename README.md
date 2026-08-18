# reflex

Small, tested helpers for Go's `reflect` package.

## Install

```text
go get github.com/Patrick-ring-motive/reflex@latest
```

```go
import "github.com/Patrick-ring-motive/reflex"
```

## Behavior

- `TypeOf`, `ValueOf`, and `AnyOf` accept ordinary values, `reflect.Value`,
  and `reflect.Type`.
- `MakeSlice` returns a concrete slice, defaults capacity to length, and
  accepts at most length and capacity.
- `SetSliceAt`, `SetField`, and `Append` accept values assignable or
  convertible to the destination type.
- `Call` treats maps, arrays, and slices as getters using its first argument
  as a key or index. Values implementing `sync.Map.Load` return its value and
  found flag. Functions join supplied argument slices, convert compatible
  values, fill missing or incompatible arguments with zero values, and ignore
  extra fixed arguments. `Call` returns `nil` rather than panicking for invalid
  functions or invocation failures.
- `FieldByName`, `MethodByName`, `MapIndex`, and `PropertyByName` return `nil`
  for a missing or inaccessible value. `PropertyByName` checks method, field,
  then map entry.
- `IsNil` recognizes untyped nil, typed nil, and invalid `reflect.Value`.
- Operations that cannot represent failure in their signature retain Go
  reflection's panic behavior for invalid indexing, dereferencing, or closing.

## Example

```go
values := reflex.MakeSlice(0, 2).([]int)
values = reflex.Append(values, int8(3)).([]int)
```

## Release

Go publishes modules from immutable semantic-version Git tags. After checks
pass on `main`, create and push a tag:

```text
git tag -a v1.0.0 -m "v1.0.0"
git push origin v1.0.0
```

Keep repository public so `proxy.golang.org` and `pkg.go.dev` can discover
releases. Publish corrections as a new version; never move or replace a tag.

## Quality checks

`bin/build.sh` verifies formatting, runs `go vet`, runs race-enabled tests,
and requires 100% statement coverage. `bin/compile.sh` compiles all packages;
`bin/lint.sh` performs non-mutating formatting and vet checks.

```text
go test -race -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```
