package example

import (
	"errors"
	"fmt"
)

func f() error {
	s, err := doSomething()
	fmt.Println(s.X) // want "field address"
	return err
}

func doSomething() (*S, error) {
	return nil, errors.New("error")
}

type S struct {
	X int
}
