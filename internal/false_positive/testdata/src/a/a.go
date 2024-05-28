package a

func f1() error {
	s, err := newS()
	if isNotNil(err) {
		return err
	}
	println(s.X) // ok because err is nil when isNotNil returns false.
	return nil
}

func f2() error {
	s, err := newS()
	if err != nil {
		s = &S{}
	}
	println(s.X) // ok because s is not nil even if err is not nil.
	return nil
}

func f3() error {
	s, err := newS()
	if err != nil && err.Error() != "expected" { // A
		return err
	}
	if err != nil && err.Error() == "expected" { // B
		return nil
	}
	println(s.X) // ok because err is nil. (A or B <=> err is nil)
	return nil
}

func isNotNil(err error) bool {
	return err != nil
}

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
