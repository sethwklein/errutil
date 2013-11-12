package errutil

import (
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestCall(t *testing.T) {
	var err error
	e := errors.New("appended error")
	AppendCall(&err, func() error {
		return e
	})
	if err != e {
		t.Errorf("%v != %v", err, e)
	}
}

type mock struct {
	w, c error
}

func OpenMock(w, c error) (io.WriteCloser, error) {
	return &mock{w, c}, nil
}

func (m *mock) Write(p []byte) (n int, err error) {
	return 0, m.w
}

func (m *mock) Close() error {
	return m.c
}

func doCall(w, c error) (err error) {
	f, err := OpenMock(w, c)
	if err != nil {
		return err
	}
	defer AppendCall(&err, f.Close)

	_, err = f.Write(nil)
	if err != nil {
		return err
	}
	return nil
}

func TestCallNilNil(t *testing.T) {
	err := doCall(nil, nil)
	if err != nil {
		t.Errorf("%v != nil", err)
	}
}

func TestCallErrNil(t *testing.T) {
	w := errors.New("write error")
	err := doCall(w, nil)
	if err != w {
		t.Errorf("%v != %v", err, w)
	}
}

func TestCallNilErr(t *testing.T) {
	c := errors.New("close error")
	err := doCall(nil, c)
	if err != c {
		t.Errorf("%v != %v", err, c)
	}
}

func TestCallErrErr(t *testing.T) {
	w := errors.New("write error")
	c := errors.New("close error")
	err := doCall(w, c)
	correct := &errorList{[]error{w, c}}
	if !reflect.DeepEqual(err, correct) {
		t.Errorf("%#v != %#v", err, correct)
	}
}
