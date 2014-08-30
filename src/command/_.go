/*
Package command provides a command-line parser with support
for sub-commands.
The commandline arguments can affect the behaviour of the
command or the application in general (e.g. logging verbosity).

Unlike the flag package, command does not stop on the first
non-flag argument. Instead it is interpreted according to the
current parsing context (as a command or command argument).
Terminating the commandline parser can be achived with "--"
in the arguments.

The following commandline notations are supported:

	* flag
	  app -verbose
	* command
	  app exec
	* flags for the application and command
	  app -verbose run -dry
	* arguments to the command
	  app test network filesystem
	* arguments to flags
	  app print -device=spool0 -priority=4
	* app run bash -- --login
	  early parsing terminator

An example making use of all features:

	app -verbose -nice=4 test ./src/test -fail -retry=2 "*.test"

Flags itself are optional, their arguments on the other hand can
be marked as mandatory (e.g. for values without safe defaults like
configuration files). Values for flags can also be read from the
environment. If not empty is satisfies the presence requirement of
a flag.

Below is an example which defines a single command, "test":

	// APPRC=.apprc app -verbose -config test "*.test" -jobs=4 -- -v
	// app -version
	// app -help
	func main() {
		var help bool = false
		var jobs int = 1
		var config *string
		var version *bool
		var verbose *int

		app := command.NewParser("app")
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
			File("")

		app.Parse()

		if help {
			app.PrintUsage()
		} else if *version {
			printVersion((*verbose) > 0)
		} else if app.Triggered(test) {
			// test.Args() = ["*.test"]
			// app.Args() = ["-v"]
			runTests(*config, jobs, test.Args(), app.Args())
		} else {
			app.WriteError(os.Stdout, "App has no default behaviour")
		}
	}

The usage message for the above example can either be written to stdout
or any other Writer instance. The output will be formatted internally.
For errors:

	app: App has no default behaviour
	Usage: app [FLAG]... [COMMAND] [FLAG]...

	Option						Meaning
	-version					Show the library version
	-help						Show the usage message
	-verbose[=LEVEL]			Set the verbosity level
								(default: 1)
	-config=FILE				Use FILE as configuration source
								(env: APPRC)

	Command						Meaning
		Option
	help						Run the test suite
		-jobs[=NUM]				Run NUM tests in parallel

	Flag processing can be terminated using --

*/
package command
