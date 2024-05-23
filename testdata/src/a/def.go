package a

func newS() (*S, error) {
	return nil, nil
}

type S struct {
	X int
}

type T struct {
	X int
}

func (t T) S() (*S, error) {
	return nil, nil
}
