package errutil

import "strings"

// ErrorList contains one or more errors.
// It's useful if you care about the errors from deferred calls
// or if you're processing several things
// and want to return more than just the first error encountered.
// The empty ErrorList is considered a pathological case
// and its behavior is undefined.
type ErrorList interface {
	// Error is here to satisfy error and for callers who don't know this
	// might be a list. It probably doesn't give ideal results.
	Error() string

	// First returns the first error in the list.
	First() error

	// Walk is how callers who know this might be a list (possibly of lists)
	// print all the items. None of the errors passed to walkFn will be
	// an ErrorList.
	Walk(walkFn func(error))
}

type errorList struct {
	a []error
}

// Firster is a subset of ErrorList used by functions that don't need the whole
// thing.
type Firster interface {
	// This is always used in such a way that Error()string will also be
	// present.
	First() error
}

// First implements ErrorList.First.
func (list *errorList) First() error {
	// They're not supposed to do this, but let's be permissive.
	if len(list.a) < 1 {
		return nil
	}
	return list.a[0]
}

// First returns ErrorList.First() if err is a Firster and err otherwise.
func First(err error) error {
	if list, ok := err.(Firster); ok {
		return list.First()
	}
	return err
}

// Walker is a subset of ErrorList used by functions that don't need the whole
// thing.
type Walker interface {
	// This is always used in such a way that Error()string will also be
	// present.
	Walk(func(error))
}

// Walk implements ErrorList.Walk.
func (list *errorList) Walk(walkFn func(error)) {
	for _, e := range list.a {
		if l, ok := e.(Walker); ok {
			l.Walk(walkFn)
		} else {
			walkFn(e)
		}
	}
}

// Walk calls ErrorList.Walk(walkFn) if err is a Walker
// and walkFn(err) otherwise.
func Walk(err error, walkFn func(error)) {
	if list, ok := err.(Walker); ok {
		list.Walk(walkFn)
		return
	}
	walkFn(err)
}

// WalkN visits the first n entries in err. It uses Walk.
func WalkN(err error, n int, walkFn func(error)) {
	type walkEnded struct{}
	fn := func(e error) {
		walkFn(e)
		n--
		if n <= 0 {
			panic(walkEnded{})
		}
	}
	defer func() {
		e := recover()
		if _, ok := e.(walkEnded); ok {
			return
		}
		panic(e)
	}()
	Walk(err, fn)
}

// Error implements ErrorList.Error.
func (list *errorList) Error() string {
	if len(list.a) == 1 {
		return list.a[0].Error()
	}
	// oh what do we tell you, simple caller?
	a := make([]string, 0, len(list.a))
	list.Walk(func(err error) {
		a = append(a, err.Error())
	})
	return strings.Join(a, "\n")
}

// Append creates an ErrorList.
// It is fine to call it with any mix of nil, error, and Walker arguments.
// It will return nil, the only error passed in, or an ErrorList as appropriate.
// If the first non-nil argument is an ErrorList returned by this function,
// it may be modified.
func Append(errs ...error) error {
	// The common case is no errors and two arguments.
	// The general loop below takes 6+ times as long to handle this.
	if len(errs) == 2 && errs[0] == nil && errs[1] == nil {
		return nil
	}

	var a []error
	var list *errorList
	for _, e := range errs {
		if e == nil {
			continue
		}
		if len(a) == 0 {
			if l, ok := e.(*errorList); ok {
				a = l.a
				list = l
				continue
			}
		}
		if l, ok := e.(Walker); ok {
			l.Walk(func(err error) {
				a = append(a, err)
			})
			continue
		}
		a = append(a, e)
	}
	switch len(a) {
	case 0:
		return nil
	case 1:
		return a[0]
	default:
		if list != nil {
			list.a = a
			return list
		}
		return &errorList{a}
	}
	panic("unreached") // for compatibility with go1
}

// AppendCall appends any error returned by the function to any existing
// errors. It is useful when returning an error from a deferred call.
//
// WARNING: WHEN DOING THAT, MAKE SURE YOU PASS A POINTER TO A NAMED RETURN
// VALUE, ELSE THE RESULT WILL BE SILENTLY DISCARDED.
//
// See silentlybroken_test.go for examples of the problem.
func AppendCall(errp *error, f func() error) {
	*errp = Append(*errp, f())
}
