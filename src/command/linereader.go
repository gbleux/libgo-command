package command

// Common functionality of command-line reader utilities
type LineReader interface {
	// Get the number of arguments
	NArg() int

	// Get the parsed arguments
	Args() []string
}
