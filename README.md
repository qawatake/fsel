# fsel

[![Go Reference](https://pkg.go.dev/badge/github.com/qawatake/fsel.svg)](https://pkg.go.dev/github.com/qawatake/fsel)
[![test](https://github.com/qawatake/fsel/actions/workflows/test.yaml/badge.svg)](https://github.com/qawatake/fsel/actions/workflows/test.yaml)

Linter `fsel` detects nil is passed to a function that does nothing for nil.

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