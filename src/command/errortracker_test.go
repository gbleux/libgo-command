package command

import (
	"testing"
)

func TestFIFO(t *testing.T) {
	e := newErrorTracker(false, true)
	first := "first"
	second := "second"

	assertSaturation(false, e, t)

	e.Store(first)

	assertSaturation(true, e, t)

	e.Store(second)

	assertSaturation(true, e, t)
	assertError(first, e, t)
}

func TestLIFO(t *testing.T) {
	e := newErrorTracker(false, false)
	first := "first"
	second := "second"

	assertSaturation(false, e, t)

	e.Store(first)

	assertSaturation(false, e, t)

	e.Store(second)

	assertSaturation(false, e, t)
	assertError(second, e, t)
}

func TestPanic(t *testing.T) {
	e := newErrorTracker(true, true)
	m := "bail"

	defer func(assertion string) {
		if r := recover(); nil == r {
			t.Error("expected panic did not occur")
		} else if assertion != r {
			t.Errorf("expected panic reason '%s' but got '%s'", assertion, r)
		}
	}(m)

	e.Store(m)
}

func TestNoOp(t *testing.T) {
	e := newErrorTracker(false, false)

	e.Store("")

	if err := e.Error(); nil != err {
		t.Error("tracker yielded an error")
	}
}

func assertSaturation(assertion bool, unit *errorTracker, t *testing.T) {
	if assertion != unit.Saturated() {
		t.Error("unexpected tracker saturation state: %v", assertion)
	}
}

func assertError(assertion string, unit *errorTracker, t *testing.T) {
	if err := unit.Error(); nil == err {
		t.Error("no error was stored")
	} else if assertion != err.Error() {
		t.Error("expected tracker error '%s' but got '%s'", assertion, err.Error())
	}
}
