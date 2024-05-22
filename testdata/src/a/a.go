package a

func f1() error {
	s, err := newS()
	println(s.X) // want "field address"
	return err
}

func f2() error {
	s, err := newS()
	if err != nil {
		return err
	}
	println(s.X) // ok because err is nil
	return nil
}

func f3() error {
	s, _ := newS()
	if s != nil {
		println(s.X) // ok because s is not nil
	}
	return nil
}

func f4() error {
	s, err := newS()
	if err != nil {
		println(s.X) // want "field address"
		return err
	}
	return nil
}

func f5() error {
	s, err := newS()
	if err != nil {
		return err
	}
	println(s.X) // ok because err is nil
	func() {
		println(err)
	}()
	return nil
}

func g1() error {
	var t T
	s, err := t.S()
	println(s.X) // want "field address"
	return err
}

func g2() error {
	var t T
	s, err := t.S()
	if err != nil {
		return err
	}
	println(s.X) // ok because err is nil
	return nil
}

func g3() error {
	var t T
	s, _ := t.S()
	if s != nil {
		println(s.X) // ok because s is not nil
	}
	return nil
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
