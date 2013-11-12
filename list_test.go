package errutil

import (
	"errors"
	"reflect"
	"testing"
)

func TestFirstOne(t *testing.T) {
	correct := errors.New("a")
	out := First(correct)
	if out != correct {
		t.Errorf("\n%#v\n!=\n%#v\n", out, correct)
	}
}

func TestFirstThree(t *testing.T) {
	correct := errors.New("a")
	e2 := errors.New("b")
	e3 := errors.New("c")
	out := First(Append(correct, e2, e3))
	if !reflect.DeepEqual(out, correct) {
		t.Errorf("\n%#v\n!=\n%#v\n", out, correct)
	}
}

func TestMulti_Merr(t *testing.T) {
	e1 := errors.New("a")
	e2 := errors.New("b")
	e3 := errors.New("c")
	e23 := Append(e2, e3)
	out := Append(e1, e23)
	correct := &errorList{[]error{e1, e2, e3}}
	if !reflect.DeepEqual(out, correct) {
		t.Errorf("\n%#v\n!=\n%#v\n", out, correct)
	}
}

func TestMulti_ErrErrX(t *testing.T) {
	e1 := errors.New("a")
	e2 := errors.New("b")
	out := Append(e2, e1)
	correct := &errorList{[]error{e2, e1}}
	if !reflect.DeepEqual(out, correct) {
		t.Errorf("\n%#v\n!=\n%#v\n", out, correct)
	}
}

func TestMulti_ErrErr(t *testing.T) {
	e1 := errors.New("a")
	e2 := errors.New("b")
	out := Append(e1, e2)
	correct := &errorList{[]error{e1, e2}}
	if !reflect.DeepEqual(out, correct) {
		t.Errorf("\n%#v\n!=\n%#v\n", out, correct)
	}
}

func TestMulti_ErrNil(t *testing.T) {
	e := errors.New("a")
	out := Append(e, nil)
	if out != e {
		t.Errorf("%#v != %#v\n", out, e)
	}
}

func TestMulti_NilErr(t *testing.T) {
	e := errors.New("a")
	out := Append(nil, e)
	if out != e {
		t.Errorf("%#v != %#v\n", out, e)
	}
}

func TestMulti_Err(t *testing.T) {
	e := errors.New("a")
	out := Append(e)
	if out != e {
		t.Errorf("%#v != %#v\n", out, e)
	}
}

func TestMulti_NilNilNil(t *testing.T) {
	out := Append(nil, nil, nil)
	if out != nil {
		t.Errorf("%#v\n", out)
	}
}

func TestMulti_NilNil(t *testing.T) {
	out := Append(nil, nil)
	if out != nil {
		t.Errorf("%#v\n", out)
	}
}

func TestMulti_Nil(t *testing.T) {
	out := Append(nil)
	if out != nil {
		t.Errorf("%#v\n", out)
	}
}
