package main

import (
	"command"
	"fmt"
	"os"
	"path"
)

func main() {
	var help bool = false
	var jobs int = 1
	var config *string
	var version *bool
	var verbose *int

	id := path.Base(os.Args[0])
	app := command.NewParser(id, true)
	test := app.Command("test", "Run the test suite")

	app.Flag("help", "Show the usage message").
		BoolVar(&help)

	test.Flag("jobs", "Run NUM tests in parallel").
		Value("NUM", false).
		IntVar(&jobs)

	version = app.Flag("version", "Show the library version").
		Bool(false)

	verbose = app.Flag("verbose", "Set the verbosity level").
		Value("LEVEL", false).
		Int(0)

	config = app.Flag("config", "Use FILE as configuration source").
		EnvironmentValue("APPRC").
		Value("FILE", true).
		File("~/.config/app/rc")

	if err := app.Parse(); nil != err {
		app.PrintError(err.Error())
	} else if help {
		app.PrintUsage()
	} else if *version {
		printVersion(id, (*verbose) > 0)
	} else if app.Triggered(test) {
		runTests(*config, jobs, test.Args(), app.Args())
	} else {
		app.PrintError("App has no default behaviour")
	}
}

func printVersion(name string, verbose bool) {
	fmt.Println("example application '" + name + "'")
	fmt.Println("libgo-command-" + command.Version)

	if verbose {
		fmt.Println("build date:", command.BuildDate)
	}
}

func runTests(config string, jobs int, glob []string, args []string) {
	fmt.Printf("exec testframework \"-c%s\" -j%d %s %s\n",
		config, jobs, args, glob)
}
