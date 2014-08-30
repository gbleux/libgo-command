package assert

import (
	"fmt"
	"reflect"
	"runtime"
)

func StringArrayEquals(r Reporter, name string, actual []string, expected []string) {
	var expect interface{}
	var result bool

	if len(actual) != len(expected) {
		assert(r, false, 1, "Array '", name, "' length does not match.\n",
			formatComparison(len(expected), len(actual), true))
		return
	}

	for i, item := range actual {
		expect = expected[i]
		result = reflect.DeepEqual(expect, item)

		assert(r, result, 1, "Array item mismatch in '"+name+"'.\n",
			formatComparison(expect, item, true))
	}
}

func Equals(r Reporter, name, actual interface{}, expected interface{}) {
	result := reflect.DeepEqual(expected, actual)

	assert(r, result, 1, name, " does not match.\n",
		formatComparison(expected, actual, true))
}

func True(r Reporter, name string, actual bool) {
	assert(r, true == actual, 1, name, " is not true")
}

func False(r Reporter, name string, actual bool) {
	assert(r, false == actual, 1, name, " is not false")
}

func That(r Reporter, actual interface{}, matcher Matcher) {
	result := matcher.Matches(actual)
	message := formatComparison(
		matcher.Describe(),
		matcher.DescribeMismatch(actual),
		false)

	assert(r, result, 1, message)
}

func ThatReason(r Reporter, reason string, actual interface{}, matcher Matcher) {
	result := matcher.Matches(actual)
	message := formatComparison(
		matcher.Describe(),
		matcher.DescribeMismatch(actual),
		false)

	assert(r, result, 1, reason+"\n"+message)
}

func formatComparison(expected interface{}, actual interface{}, prefixActual bool) string {
	var prefix string = ""

	if prefixActual {
		prefix = "was "
	}

	return fmt.Sprintf("Expected: %v\n     but: %s%v", expected, prefix, actual)
}

func assert(r Reporter, result bool, stackPos int, args ...interface{}) {
	if false == result {
		_, file, line, _ := runtime.Caller(stackPos + 1)
		message := fmt.Sprint(args...)

		r.Errorf("%s\n(%s:%d)", message, file, line)
	}
}
