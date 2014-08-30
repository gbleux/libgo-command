package command

import (
	"testing"
)

type booleanTest struct {
	input    string
	expected bool
}

func TestBooleanValue(t *testing.T) {
	var success = []booleanTest{
		{"true", true},
		{"True", true},
		{"yes", true},
		{"1", true},
		{"", true},
		{"no", false},
	}
	var failure = []string{
		"tRue",
		"nope",
		"-",
	}
	var actual bool
	var unit boolValue

	if false == unit.IsBoolFlag() {
		t.Error("boolean value is not identifying itself as boolean flag")
	}

	for _, pair := range success {
		actual = false
		unit = boolValue{&actual}

		if err := unit.Set(pair.input); nil != err {
			t.Error("setting boolean value to", pair.input, "yielded an error:", err.Error())
		} else if actual != pair.expected {
			t.Error("boolean mismatch. expected: ", pair.expected, "got:", actual)
		}
	}

	for _, input := range failure {
		actual = false
		unit = boolValue{&actual}

		if err := unit.Set(input); nil == err {
			t.Error("invalid boolean input", input, "caused no error")
		}
	}
}
