// This example provides insight into its own logic and the
// internals of libgo-command.
package main

import (
	"fmt"

	"command"
)

func main() {
	var shared int = 1

	app := command.NewParser("playground", true)

	app.
		Flag("shared1", "set the value of the shared variable").
		Value("NUMBER", true).
		IntVar(&shared)
	app.
		Flag("shared2", "set the value of the shared variable").
		Value("INTEGER", false).
		IntVar(&shared)
	app.
		Flag("deprecated", "ignored").
		Void(false)
	app.
		Flag("legacy", "ignored").
		Void(true)

	if err := app.Parse(); nil != err {
		app.PrintError(err.Error())
		return
	}

	fmt.Println("Shared value is", shared)
}
