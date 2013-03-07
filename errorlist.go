package errors

import "strings"

// ErrorList contains one or more errors.
// It's useful if you care about the errors from deferred calls
// or if you're processing several things
// and want to return more than just the first error encountered.
// The empty ErrorList is considered a pathological case
// and its behavior is undefined.
type ErrorList interface {
	// Error is here to satisfy error and for callers who don't know this
	// might be a list. It probably doesn't give very good results.
	Error() string

	// Walk is how callers who know this might be a list (possibly of lists)
	// print all the items. None of the errors passed to walkFn will be
	// an ErrorList.
	Walk(walkFn func(error))
}

type errorList struct {
	a []error
}

// Walker is a subset of ErrorList used by functions that require only it.
type Walker interface {
	Walk(func(error))
}

// Walk calls ErrorList.Walk if err is a Walker and walkFn(err) otherwise.
func Walk(err error, walkFn func(error)) {
	if list, ok := err.(Walker); ok {
		list.Walk(walkFn)
		return
	}
	walkFn(err)
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
	panic("unreached")
}
