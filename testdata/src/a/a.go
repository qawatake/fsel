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
	func() { println(err) }()
	return nil
}

func f6() error {
	s, err := newS()
	if err != nil {
		return err
	}
	s, err = newS()
	println(s.X) // want "field address"
	func() { println(err) }()
	return nil
}

func f7() error {
	s, err := newS()
	if err != nil {
		return err
	}
	s, err = newS()
	if err != nil {
		return err
	}
	println(s.X) // ok
	func() { println(err) }()
	return nil
}

func f8() (err error) {
	s, err := newS()
	if err != nil {
		return err
	}
	println(s.X) // ok because err is nil
	return nil
}

func f9() (err error) {
	s, err := newS()
	println(s.X) // want "field address"
	return err
}

func f10() error {
	s, err := newS()
	println(s.X) // want "field address"
	func() { println(s.X) }()
	return err
}

func f11() (err error) {
	s, err := newS()
	if s != nil {
		if true {
			return nil
		}
	}
	if err != nil {
		return err
	}
	println(s.X) // ok because err is nil
	func() { println(err) }()
	return nil
}

func f12() error {
	s, err := newS()
	if err != nil {
		println(s.X) // want "field address"
		return err
	}
	func() { println(s.X) }()
	return nil
}

func f13() error {
	s, err := newS()
	println(s.X) //lint:ignore fsel reason
	return err
}

func f14() error {
	s, err := newS()
	if s.X == 0 { //lint:ignore fsel reason
		return err
	}
	return nil
}

func f15() error {
	s, err := newS()
	//lint:ignore fsel reason
	println(s.X)
	return err
}

func f16() error {
	s, err := newS()
	//lint:ignore fsel reason
	if s.X == 0 {
		return err
	}
	return nil
}

func f17() error {
	s, err := newS()
	//lint:ignore fsel reason
	println(s.X)
	println(s.X)
	return err
}

func f18() error {
	s, err := newS()
	//lint:ignore fsel reason
	println(s.X)
	println(s.X)
	defer func() { println(s) }()
	return err
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
