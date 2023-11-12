package example

import (
	"errors"
	"fmt"
)

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
