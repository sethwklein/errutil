package errutil

import (
	"errors"
	"reflect"
	"testing"
)

func TestInvalidNil(t *testing.T) {
	got := (&errorList{}).First()
	if got != nil {
		t.Errorf("expected nil, got: %v\n", got)
	}
}

func TestInvalidLenOne(t *testing.T) {
	want := "want"
	list := &errorList{a: []error{errors.New(want)}}
	got := list.Error()
	if want != got {
		t.Errorf("wanted: %v, got: %v\n", want, got)
	}
}

func TestNested(t *testing.T) {
	chick := errors.New("in nest")
	nest := &errorList{a: []error{
		&errorList{a: []error{ chick }},
	}}
	var got error
	nest.Walk(func(err error) {
		got = err
	})
	if got != chick {
		t.Errorf("expected: %v, got: %v\n", chick, got)
	}
}

func TestWalkNil(t *testing.T) {
	called := false
	Walk(nil, func(err error) {
		called = true
	})
	if called {
		t.Error("func was called but should not have been")
	}
}

func TestWalkNNil(t *testing.T) {
	called := false
	WalkN(nil, 1, func(err error) {
		called = true
	})
	if called {
		t.Error("func was called but should not have been")
	}
}

func TestWalkSingle(t *testing.T) {
	want := errors.New("want")
	var got error
	Walk(want, func(err error) {
		got = err
	})
	if want != got {
		t.Errorf("wanted: %v, got: %v\n", want, got)
	}
}

func TestWalkNPanic(t *testing.T) {
	want := errors.New("want")
	defer func() {
		got := recover()
		if got != want {
			t.Errorf("wanted: %v, got: %v\n", want, got)
		}
	}()
	WalkN(errors.New("dummy"), 1, func(_ error) {
		panic(want)
	})
}

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
