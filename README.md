# fsel

[![Go Reference](https://pkg.go.dev/badge/github.com/qawatake/fsel.svg)](https://pkg.go.dev/github.com/qawatake/fsel)
[![test](https://github.com/qawatake/fsel/actions/workflows/test.yaml/badge.svg)](https://github.com/qawatake/fsel/actions/workflows/test.yaml)

Linter: `fsel` flags field access with unverified nil errors.

```go
func bad() error {
  s, err := doSomething()
  fmt.Println(s.X) // <- field address without checking nilness of err
  return err
}

func good() error {
  s, err := doSomething()
  if err != nil {
    return err
  }
  fmt.Println(s.X) // ok because err is definitely nil
  return nil
}

func doSomething() (*S, error) {
  return nil, errors.New("error")
}

type S struct {
  X int
}
```

You can try an example by running `make run.example`.

## How to use

```sh
go install github.com/qawatake/fsel/cmd/fsel@latest
fsel ./...
```

## Dealing with False Positives

In case of a false positive, consider either adding an ignore comment or refactoring your code.
The latter is preferable as it enhances code resilience to future changes.

Example of a false positive:

```go
func f() error {
  s, err := doSomething()
  if isNotNil(err) {
    return err
  }
  fmt.Println(s.X) // <- false positive??
  return nil
}
```

Refactoring example:

```go
func f() error {
  s, err := doSomething()
  if isNotNil(err) {
    return err
  }
  if s == nil {
    return errors.New("unreachable")
  }
  fmt.Println(s.X) // ok
  return nil
}
```

Example of using an ignore comment:

```go
func f() error {
  s, err := doSomething()
  if isNotNil(err) {
    return err
  }
  fmt.Println(s.X) // lint:ignore fsel reason
  return nil
}
```
