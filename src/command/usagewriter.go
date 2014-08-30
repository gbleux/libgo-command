package command

import (
	"fmt"
	"io"
)

const (
	formatUsage        = "[FLAG]... [COMMAND] [FLAG]..."
	formatColumn       = "%-31s %s\n"
	formatIndent       = "    "
	formatFlagBool     = flagPrefix + "%s"
	formatFlagRequired = flagPrefix + "%s" + flagValueSep + "%s"
	formatFlagOptional = flagPrefix + "%s" + flagValueSep + "[%s]"
)

type usageWriter struct {
	out io.Writer
}

// Write the command usage pattern.
func (u *usageWriter) WriteTitle(application string) *usageWriter {
	fmt.Fprintln(u.out, "Usage:", application, formatUsage)

	return u
}

// Describe a single commandline flag.
func (u *usageWriter) WriteFlag(name string, flag *Flag, indent bool) *usageWriter {
	prefix := formatFlag(name, flag.valueName, flag.valueReq, flag.value)

	if indent {
		prefix = formatIndent + prefix
	}

	fmt.Fprintf(u.out, formatColumn, prefix, flag.desc)

	return u
}

// Describe a command.
func (u *usageWriter) WriteCommand(name string, cmd *Command) *usageWriter {
	fmt.Fprintf(u.out, formatColumn, name, cmd.desc)

	return u
}

// Write the command/flag description header. The command header
// has three columns streched over two lines, whereas the flag header
// uses only a single line.
func (u *usageWriter) WriteHeader(commandHeader bool) *usageWriter {
	fmt.Fprintln(u.out)

	if commandHeader {
		fmt.Fprintf(u.out, formatColumn, "Command", "Meaning")
		fmt.Fprintln(u.out, formatIndent+"Option")
	} else {
		fmt.Fprintf(u.out, formatColumn, "Option", "Meaning")
	}

	return u
}

// Write the usage footer message.
func (u *usageWriter) WriteFooter() *usageWriter {
	fmt.Fprintln(u.out)
	fmt.Fprintln(u.out, "Flag processing can be terminated using", flagTermination)

	return u
}

func newUsageWriter(writer io.Writer) *usageWriter {
	return &usageWriter{writer}
}

func formatFlag(name string, key string, req bool, value interface{}) string {
	if b, ok := value.(inferableValue); ok && b.IsBoolFlag() {
		return fmt.Sprintf(formatFlagBool, name)
	} else if req {
		return fmt.Sprintf(formatFlagRequired, name, key)
	}

	return fmt.Sprintf(formatFlagOptional, name, key)
}
