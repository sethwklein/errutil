package errors_test

import (
	"fmt"
	"sethwklein.net/errors"
	"testing"
)

const listLen = 40
const attentionSpan = 30

func BenchmarkWalkPanic(b *testing.B) {
	var list error
	for i := 1; i <= listLen; i++ {
		list = errors.Append(list, fmt.Errorf("number %v", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				if e := recover(); e != nil {
					panic(e)
				}
			}()
			n := attentionSpan
			i := 0
			errors.Walk(list, func(e error) {
				_ = e
				i++
				if i >= n {
					panic(nil)
				}
			})
		}()
	}
}

func BenchmarkWalkIgnore(b *testing.B) {
	var list error
	for i := 1; i <= listLen; i++ {
		list = errors.Append(list, fmt.Errorf("number %v", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := attentionSpan
		i := 0
		errors.Walk(list, func(e error) {
			i++
			if i > n {
				return
			}
			_ = e
		})
	}
}

func BenchmarkWalkN(b *testing.B) {
	var list error
	for i := 1; i <= listLen; i++ {
		list = errors.Append(list, fmt.Errorf("number %v", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		errors.WalkN(list, attentionSpan, func(e error) {
			_ = e
		})
	}
}
