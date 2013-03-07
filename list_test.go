package errors

import (
	"reflect"
	"testing"
)

func TestFirstOne(t *testing.T) {
	correct := New("a")
	out := First(correct)
	if out != correct {
		t.Logf("\n%#v\n!=\n%#v\n", out, correct)
		t.Fail()
	}
}

func TestFirstThree(t *testing.T) {
	correct := New("a")
	e2 := New("b")
	e3 := New("c")
	out := First(Append(correct, e2, e3))
	if !reflect.DeepEqual(out, correct) {
		t.Logf("\n%#v\n!=\n%#v\n", out, correct)
		t.Fail()
	}
}

func TestMulti_Merr(t *testing.T) {
	e1 := New("a")
	e2 := New("b")
	e3 := New("c")
	e23 := Append(e2, e3)
	out := Append(e1, e23)
	correct := &errorList{[]error{e1, e2, e3}}
	if !reflect.DeepEqual(out, correct) {
		t.Logf("\n%#v\n!=\n%#v\n", out, correct)
		t.Fail()
	}
}

func TestMulti_ErrErrX(t *testing.T) {
	e1 := New("a")
	e2 := New("b")
	out := Append(e2, e1)
	correct := &errorList{[]error{e2, e1}}
	if !reflect.DeepEqual(out, correct) {
		t.Logf("\n%#v\n!=\n%#v\n", out, correct)
		t.Fail()
	}
}

func TestMulti_ErrErr(t *testing.T) {
	e1 := New("a")
	e2 := New("b")
	out := Append(e1, e2)
	correct := &errorList{[]error{e1, e2}}
	if !reflect.DeepEqual(out, correct) {
		t.Logf("\n%#v\n!=\n%#v\n", out, correct)
		t.Fail()
	}
}

func TestMulti_ErrNil(t *testing.T) {
	e := New("a")
	out := Append(e, nil)
	if out != e {
		t.Logf("%#v != %#v\n", out, e)
		t.Fail()
	}
}

func TestMulti_NilErr(t *testing.T) {
	e := New("a")
	out := Append(nil, e)
	if out != e {
		t.Logf("%#v != %#v\n", out, e)
		t.Fail()
	}
}

func TestMulti_Err(t *testing.T) {
	e := New("a")
	out := Append(e)
	if out != e {
		t.Logf("%#v != %#v\n", out, e)
		t.Fail()
	}
}

func TestMulti_NilNilNil(t *testing.T) {
	out := Append(nil, nil, nil)
	if out != nil {
		t.Logf("%#v\n", out)
		t.Fail()
	}
}

func TestMulti_NilNil(t *testing.T) {
	out := Append(nil, nil)
	if out != nil {
		t.Logf("%#v\n", out)
		t.Fail()
	}
}

func TestMulti_Nil(t *testing.T) {
	out := Append(nil)
	if out != nil {
		t.Logf("%#v\n", out)
		t.Fail()
	}
}
