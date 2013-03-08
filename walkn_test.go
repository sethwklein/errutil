package errors_test

import (
	"fmt"
	"os"
	"sethwklein.net/errors"
)

func ExampleWalkN() {
	var list error
	for i := 1; i <= 1000; i++ {
		list = errors.Append(list, fmt.Errorf("number %v", i))
	}
	errors.WalkN(list, 3, func(e error) {
		fmt.Fprintln(os.Stderr, e)
	})
	// Output:
	// number 1
	// number 2
	// number 3
}
