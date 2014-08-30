package command

type Command struct {
	// dynamic initialization data

	flags map[string]*Flag
	desc  string

	// dynamic runtime data

	args []string
}

func (c *Command) Flag(name string, description string) *Flag {
	flag := newFlag(description)

	c.flags[name] = flag

	return flag
}

// Returns the number of arguments remaining
// after flags have been processed
func (c *Command) NArg() int {
	return len(c.args)
}

func (c *Command) Args() []string {
	return c.args
}

func (c *Command) String() string {
	return c.desc
}

func newCommand(description string) *Command {
	flags := make(map[string]*Flag)
	args := []string{}

	return &Command{flags, description, args}
}
