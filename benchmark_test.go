package errutil_test

import (
	"fmt"
	"sethwklein.net/go/errutil"
	"testing"
)

func BenchmarkBase(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		err = fmt.Errorf("number %v", i)
	}
	_ = err
}

func BenchmarkAppend(b *testing.B) {
	var list error
	for i := 0; i < b.N; i++ {
		list = errutil.Append(list, fmt.Errorf("number %v", i))
	}
}

func BenchmarkNil(b *testing.B) {
	var list error
	for i := 0; i < b.N; i++ {
		list = errutil.Append(list, nil)
	}
}
