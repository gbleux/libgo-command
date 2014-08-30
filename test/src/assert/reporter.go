package assert

// Assertion status reporter. Emits messages about failed assertions
// as well as general test progression.
// This interface is intended to avoid hard dependencies on
// [testing.T].
type Reporter interface {
	// Report a failed assertion
	Error(args ...interface{})
	// Report a failed assertion
	Errorf(format string, args ...interface{})
	// Report a message without failing the assertion.
	Log(args ...interface{})
	// Report a message without failing the assertion.
	Logf(format string, args ...interface{})
}
