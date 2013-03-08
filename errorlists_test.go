// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Portions copyright 2013 Seth W. Klein.

package errors_test

import (
	"sethwklein.net/go/errors"
	"testing"
)

func TestNewEqualList(t *testing.T) {
	// Different allocations should not be equal.
	if errors.Append(errors.New("ab"), errors.New("cd")) == errors.Append(errors.New("ab"), errors.New("cd")) {
		t.Errorf(`Append(New("ab"), New("cd")) == Append(New("ab"), New("cd"))`)
	}

	// Same allocation should be equal to itself (not crash).
	err := errors.Append(errors.New("jk"), errors.New("lm"))
	if err != err {
		t.Errorf(`err != err`)
	}
}

func TestErrorMethodList(t *testing.T) {
	err := errors.Append(errors.New("ab"), errors.New("cd"))
	if err.Error() != "ab"+"\n"+"cd" {
		t.Errorf(`Append(New("ab"), New("cd")) = %q, want %q`,
			err.Error(), "ab\ncd")
	}
}
