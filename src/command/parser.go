package command

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	flagTermination = "--"
	flagPrefix      = "-"
	flagValueSep    = "="
)

type Parser struct {
	// static data

	lenient bool
	owner   string

	// dynamic initialization data

	flags map[string]*Flag
	cmds  map[string]*Command

	// dynamic runtime data

	args    []string
	trigger *Command
	invokes int
}

// Register a new command name. The description is included in the
// auto-generated usage message along with the name. The name is
// case sensitive.
func (p *Parser) Command(name string, description string) *Command {
	cmd := newCommand(description)

	p.cmds[name] = cmd

	return cmd
}

// Register a new flag to alter the behaviour of the application.
// The flag is case sensitive.
func (p *Parser) Flag(name string, description string) *Flag {
	flag := newFlag(description)

	p.flags[name] = flag

	return flag
}

func (p *Parser) NArg() int {
	return len(p.args)
}

func (p *Parser) Args() []string {
	return p.args
}

func (p *Parser) WriteError(out io.Writer, message string) {
	fmt.Fprintf(out, "%s: %s\n", p.owner, message)

	p.WriteUsage(out)
}

func (p *Parser) WriteUsage(out io.Writer) {
	usage := newUsageWriter(out)

	usage.WriteTitle(p.owner)

	if len(p.flags) > 0 {
		usage.WriteHeader(false)
	}

	for name, flag := range p.flags {
		usage.WriteFlag(name, flag, false)
	}

	if len(p.cmds) > 0 {
		usage.WriteHeader(true)
	}

	for name, cmd := range p.cmds {
		usage.WriteCommand(name, cmd)

		for key, flag := range cmd.flags {
			usage.WriteFlag(key, flag, true)
		}
	}

	usage.WriteFooter()
}

// Proxy method for WriteError using os.Stderr as output writer.
func (p *Parser) PrintError(message string) {
	p.WriteError(os.Stderr, message)
}

// Proxy method for WriteUsage using os.Stdout as output writer.
func (p *Parser) PrintUsage() {
	p.WriteUsage(os.Stdout)
}

// Process the provided arguments. The slice is expected to contain
// only flags, commands and command flags, i.e. the application name
// from os.Args or similar slices should be omitted.
// The return value indicates problems during processing. If the
// parser is configured to continue on errors, the last encountered
// error is returned.
// Calling this method multiple times will overwrite any previous
// parsing results as it is reset first before writing new data.
// (see Reset())
func (p *Parser) ParseArgs(argv []string) (err error) {
	// use root flags first
	var flags map[string]*Flag = p.flags
	var cmd *Command = nil
	var e *errorTracker = newErrorTracker(!p.lenient, true)
	var args []string = make([]string, 0, len(argv)) // pessimistic size
	var cmdArgs bool = false                         // where to append non-flags

	defer func() {
		if r := recover(); nil != r {
			// the error might also come from areas
			// not covered by the error tracker
			err = fmt.Errorf("%v", r)
		}

		if nil == cmd {
			cmd = newCommand("") // dummy value
		}

		// even though we might have encountered
		// errors, the parser is configured to
		// complete the parsing anyway. the following
		// values contain values which where successfully
		// parsed from the arguments
		p.trigger = cmd
		cmd.args = args
	}()

	p.args = make([]string, 0, 0) // will stay this way or overwritten
	p.invokes++

	for index, arg := range argv {
		if arg == flagTermination {
			// we do not want the flag terminator in the array
			if len(argv) > index {
				p.args = argv[index+1:]
			}
			break
		} else if strings.HasPrefix(arg, flagPrefix) {
			e.Store(parseFlag(arg, flags))
		} else if cmdArgs {
			// append to command args
			args = append(args, arg)
		} else if cmd, cmdArgs = p.cmds[arg]; cmdArgs {
			// use command flags from now on.
			flags = cmd.flags
		} else {
			// argument is neither a flag nor a valid command
			e.Store(fmt.Sprintf("No such command '%s'", arg))
		}
	}

	return e.Error()
}

// This method is similar to the Parse method of the flag package.
// It is a simple proxy method calling ParseArgs with os.Args[1:].
func (p *Parser) Parse() error {
	return p.ParseArgs(os.Args[1:])
}

// Returns true if Parse has been called at least once.
func (p *Parser) Parsed() bool {
	return p.invokes > 0
}

// Returns true if the provided command was triggered during the
// last parsing process.
func (p *Parser) Triggered(cmd *Command) bool {
	return p.trigger == cmd
}

// Create a new parser instance. The application name is provided
// here as it is not expected to be part of any parsing input. It is
// used in the usage message.
// If an error occurs during the processing, the parser can bail out
// immediately or continue with the next argument, ignoring any invalid
// data.
func NewParser(applicationName string, continueOnError bool) *Parser {
	flags := make(map[string]*Flag)
	cmds := make(map[string]*Command)
	args := []string{}

	return &Parser{continueOnError,
		applicationName,
		flags,
		cmds,
		args,
		nil,
		0}
}

func parseFlag(needle string, haystack map[string]*Flag) string {
	flag := strings.TrimPrefix(needle, flagPrefix)
	parts := strings.SplitN(flag, flagValueSep, 2)
	key := parts[0]
	val := ""

	if len(parts) > 1 {
		val = parts[1]
	}

	if flag, ok := haystack[key]; ok {
		if err := flag.value.Set(val); nil != err {
			if msg := err.Error(); len(msg) > 0 {
				return msg
			} else {
				return fmt.Sprintf("Unable to write value '%s' to flag '%s'", val, key)
			}
		}

		return ""
	}

	return fmt.Sprintf("No such flag '%s'", key)
}
