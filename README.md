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

## False Positives

To ignore a false positive, add a comment `//lint:ignore fsel reason` to the line.

```go
func f() error {
  s, err := doSomething()
  if isNotNil(err) {
    return err
  }
  fmt.Println(s.X) //lint:ignore fsel reason
  return nil
}
```
