package a

import (
	"errors"
	"fmt"
)

func f() {
	s, err := doSomething()
	fmt.Println(s.X) // want "s.X is dereferenced without checking that it is not nil"
}

func doSomething() (*S, error) {
	return nil, errors.New("error")
}

type S struct {
	X int
}
