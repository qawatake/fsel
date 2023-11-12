package a

func f() error {
	s, err := doSomething()
	println(s.X) // want "field address"
	return err
}

func doSomething() (*S, error) {
	// return nil, errors.New("error")
	return nil, nil
}

type S struct {
	X int
}
