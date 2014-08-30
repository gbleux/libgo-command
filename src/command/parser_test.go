package command

import (
	"fmt"
	"testing"

	"assert"
)

type parserMatcher struct {
	parser  *Parser
	cmdline []string
	err     string
}

func (p *parserMatcher) Matches(item interface{}) bool {
	if err := p.parser.ParseArgs(p.cmdline); nil != err {
		p.err = fmt.Sprint("failed to parse arguments:", err)

		return false
	}

	if command, ok := item.(*Command); ok {
		if false == p.parser.Triggered(command) {
			p.err = fmt.Sprint("expected command", command, "was not triggered")
			return false
		}

		return true
	}

	return false
}

func (p *parserMatcher) DescribeMismatch(item interface{}) string {
	if 0 < len(p.err) {
		return p.err
	}

	return fmt.Sprint("unsuccessful parsing")
}

func (p *parserMatcher) Describe() string {
	return "successful argument parsing with expected command"
}

// aflags, command, cflags, cargs -- aargs
func TestParseArgsFull(t *testing.T) {
	a := false
	c := false
	aargs := []string{"app_arg1", "app_arg2"}
	cargs := []string{"cmd_arg1", "cmd_arg2", "cmd_arg3"}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(true, true, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.True(t, "app flag", a)
	assert.True(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 2)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 3)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// command, cflags, cargs -- aargs
func TestParseArgsCcFcAaA(t *testing.T) {
	a := false
	c := false
	aargs := []string{"app_arg1", "app_arg2"}
	cargs := []string{"cmd_arg1", "cmd_arg2", "cmd_arg3"}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(false, true, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.False(t, "app flag", a)
	assert.True(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 2)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 3)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// aflags, command, cargs -- aargs
func TestParseArgsaFCcAaA(t *testing.T) {
	a := false
	c := false
	aargs := []string{"app_arg1", "app_arg2"}
	cargs := []string{"cmd_arg1", "cmd_arg2", "cmd_arg3"}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(true, false, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.True(t, "app flag", a)
	assert.False(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 2)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 3)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// aflags, command, cflags, -- aargs
func TestParseArgsaFCcFaA(t *testing.T) {
	a := false
	c := false
	aargs := []string{"app_arg1", "app_arg2"}
	cargs := []string{}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(true, true, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.True(t, "app flag", a)
	assert.True(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 2)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 0)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// aflags, command, cflags, cargs
func TestParseArgsaFCcFcA(t *testing.T) {
	a := false
	c := false
	aargs := []string{}
	cargs := []string{"cmd_arg1", "cmd_arg2", "cmd_arg3"}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(true, true, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.True(t, "app flag", a)
	assert.True(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 0)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 3)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// command, cargs -- aargs
func TestParseArgsCcAaA(t *testing.T) {
	a := false
	c := false
	aargs := []string{"app_arg1", "app_arg2"}
	cargs := []string{"cmd_arg1", "cmd_arg2", "cmd_arg3"}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(false, false, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.False(t, "app flag", a)
	assert.False(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 2)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 3)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// command, -- aargs
func TestParseArgsCaA(t *testing.T) {
	a := false
	c := false
	aargs := []string{"app_arg1", "app_arg2"}
	cargs := []string{}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(false, false, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.False(t, "app flag", a)
	assert.False(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 2)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 0)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// command, cargs
func TestParseArgsCcA(t *testing.T) {
	a := false
	c := false
	aargs := []string{}
	cargs := []string{"cmd_arg1", "cmd_arg2", "cmd_arg3"}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(false, false, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.False(t, "app flag", a)
	assert.False(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 0)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 3)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// command, cflags
func TestParseArgsCcF(t *testing.T) {
	a := false
	c := false
	aargs := []string{}
	cargs := []string{}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(false, true, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.False(t, "app flag", a)
	assert.True(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 0)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 0)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

// command
func TestParseArgsC(t *testing.T) {
	a := false
	c := false
	aargs := []string{}
	cargs := []string{}
	unit, cmd := newParserTestUnit(&a, &c)
	arguments, _ := newArgs(false, false, aargs, cargs)

	assert.That(t, cmd, newParserMatcher(unit, arguments))

	assert.False(t, "app flag", a)
	assert.False(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 0)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 0)

	assert.StringArrayEquals(t, "app args", unit.Args(), aargs)
	assert.StringArrayEquals(t, "cmd args", cmd.Args(), cargs)
}

func TestParseArgsNotTriggered(t *testing.T) {
	a := false
	c := false
	unit, cmd := newParserTestUnit(&a, &c)
	argv := []string{"cmd2", "-dflag2"}

	if err := unit.ParseArgs(argv); nil != err {
		t.Fatal("failed to parse arguments")
	}

	if unit.Triggered(cmd) {
		t.Fatal("cmd1 was falsely marked as triggered")
	}
}

func TestParseArgsAppFlagsOnly(t *testing.T) {
	a := false
	c := false
	unit, cmd := newParserTestUnit(&a, &c)
	argv := []string{"-aflag1"}

	if err := unit.ParseArgs(argv); nil != err {
		t.Fatal("failed to parse arguments")
	}

	assert.True(t, "app flag", a)
	assert.False(t, "cmd flag", c)

	assert.Equals(t, "app arg count", unit.NArg(), 0)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 0)
}

func TestParseArgsAppArgsOnly(t *testing.T) {
	a := false
	c := false
	unit, cmd := newParserTestUnit(&a, &c)
	argv := []string{"--", "--login"}

	if err := unit.ParseArgs(argv); nil != err {
		t.Fatal("failed to parse arguments")
	}

	assert.Equals(t, "app arg count", unit.NArg(), 1)
	assert.Equals(t, "cmd arg count", cmd.NArg(), 0)

	assert.StringArrayEquals(t, "app args", unit.Args(), argv[1:])
}

func newArgs(appFlag bool, cmdFlag bool, appArgs []string, cmdArgs []string) ([]string, int) {
	var aa int = len(appArgs)
	var ca int = len(cmdArgs)
	var size int = 1 + 1 + aa + ca + 1 + 1
	var array []string = make([]string, 0, size)

	if appFlag {
		array = append(array, "-aflag1")
	}

	array = append(array, "cmd1")

	if cmdFlag {
		array = append(array, "-cflag1")
	}

	array = append(array, cmdArgs...)

	if 0 < len(appArgs) {
		array = append(array, "--")
	}

	array = append(array, appArgs...)

	return array, len(array)
}

func newParserTestUnit(parseFlag *bool, cmdFlag *bool) (*Parser, *Command) {
	parser := NewParser("testing", true)
	command1 := parser.Command("cmd1", "test command 1")
	command2 := parser.Command("cmd2", "test command 2")

	parser.Flag("aflag1", "test flag").BoolVar(parseFlag)
	command1.Flag("cflag1", "test flag").BoolVar(cmdFlag)
	command2.Flag("dflag1", "test flag").Bool(false)
	command2.Flag("dflag2", "test flag").Bool(false)

	return parser, command1
}

func newParserMatcher(p *Parser, argv []string) *parserMatcher {
	return &parserMatcher{p, argv, ""}
}
