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

	// WalkPartial is like Walk, but if walkFn returns false,
	// processing stops. The return value is true if walkFn never
	// returned false.
	WalkPartial(walkFn func(error) bool) bool
}

type errorList struct {
	a []error
}

// Walk calls ErrorList.Walk if err is an ErrorList and walkFn(err) otherwise.
func Walk(err error, walkFn func(error)) {
	if list, ok := err.(ErrorList); ok {
		list.Walk(walkFn)
		return
	}
	walkFn(err)
}

func (list *errorList) Walk(walkFn func(error)) {
	for _, e := range list.a {
		if l, ok := e.(ErrorList); ok {
			l.Walk(walkFn)
		} else {
			walkFn(e)
		}
	}
}

// WalkPartial calls ErrorList.WalkPartial if err is an ErrorList and
// walkFn(err) otherwise.
func WalkPartial(err error, walkFn func(error) bool) bool {
	if list, ok := err.(ErrorList); ok {
		return list.WalkPartial(walkFn)
	}
	return walkFn(err)
}

func (list *errorList) WalkPartial(walkFn func(error) bool) bool {
	again := true
	for _, e := range list.a {
		if l, ok := e.(ErrorList); ok {
			again = l.WalkPartial(walkFn)
		} else {
			again = walkFn(e)
		}
		if !again {
			break
		}
	}
	return again
}

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
// It is fine to call it with any mix of nil, error, and ErrorList arguments.
// It will return nil, the only error passed in, or an ErrorList as appropriate.
func Append(errs ...error) error {
	var a []error
	for _, e := range errs {
		if e == nil {
			continue
		}
		if l, ok := e.(ErrorList); ok {
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
		return &errorList{a}
	}
	panic("unreached")
}
