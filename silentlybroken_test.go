package errors_test

// This file demonstrates how if you forget to name your return value, no
// matter how you handle errors from deferred calls, you have a bug only
// advanced testing will point out.
//
// All the broken* functions below should return the error from mockFile.Close,
// but they don't because they forget to use a named return value.

import (
	"io"
	"sethwklein.net/go/errors"
	"testing"
)

var mockOS = struct {
	Create func(string) (io.WriteCloser, error)
}{
	func(_ string) (io.WriteCloser, error) {
		return mockFile{}, nil
	},
}

type mockFile struct{}

var mockError = errors.New("mock error")

func (_ mockFile) Close() error {
	return mockError
}

func (_ mockFile) Write(data []byte) (n int, err error) {
	return len(data), nil
}

func brokenAppendWriteFile(filename string, data []byte) error {
	f, err := mockOS.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Append(err, f.Close())
	}()
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return err
}

func brokenCallWriteFile(filename string, data []byte) error {
	f, err := mockOS.Create(filename)
	if err != nil {
		return err
	}
	defer errors.AppendCall(&err, f.Close)
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return err
}

func brokenManualWriteFile(filename string, data []byte) error {
	f, err := mockOS.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		if e != nil && err == nil {
			// This is assumes both errors can't be meaningful
			// at the same time which is not guaranteed by the
			// method signatures or documentation, although it
			// may be guaranteed by the implementation.
			err = e
		}
	}()
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	return err
}

func TestBrokenAppend(t *testing.T) {
	err := brokenAppendWriteFile("example.txt", []byte("example!"))
	// expected fail
	if err != nil { // err != mockError {
		t.Errorf("%v != %v", err, mockError)
	}
}

func TestBrokenCall(t *testing.T) {
	err := brokenCallWriteFile("example.txt", []byte("example!"))
	// expected fail
	if err != nil { // err != mockError {
		t.Errorf("%v != %v", err, mockError)
	}
}

func TestBrokenManual(t *testing.T) {
	err := brokenManualWriteFile("example.txt", []byte("example!"))
	// expected fail
	if err != nil { // err != mockError {
		t.Errorf("%v != %v", err, mockError)
	}
}
