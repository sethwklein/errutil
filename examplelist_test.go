package errutil_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sethwklein.net/go/errutil"
)

func ReinventTheIOUtil(filename string) (buf []byte, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer errutil.AppendCall(&err, f.Close)

	buf, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ExampleAppendCall() {
	buf, err := ReinventTheIOUtil("errorlist.go")
	if err != nil {
		command := filepath.Base(os.Args[0])
		errutil.Walk(err, func(e error) {
			fmt.Fprintf(os.Stderr, "%s, Error: %s\n", command, e)
		})
	}
	_ = buf // do something magical with it!
}
