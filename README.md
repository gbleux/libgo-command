# libgo-command

A command-line parser library for programs using sub-commands.

## Install

Run the command

> go install github.com/gbleux/libgo-command

_Go_ will probably emit a warning about missing compile units.
The repository will end up in **$GOPATH** and can be imported
with

    import (
        "github.com/gbleux/libgo-command/src/command"
    )

## Usage

Example applications using _libgo-command_ can be found in the
**src/examples** directory.

### Compatibility with _flag_

While _libgo-command_ is not a direct drop-in replacement for the
_flag_ package, it inherits some concepts. The main difference is
the creation of the command-line specification.

Compare a simple _flag_-base example ...

    package main
    
    import "flag"
    
    func main() {
        number := flag.Int("num", 1, "flag description")
    
        flag.Parse()
    
        // application logic ...
    }

with the same logic implemented using _libgo-command_:

    package main
    
    import "command"
    
    func main() {
        app := command.NewParser("example", true)
        number := app.Flag("num", "flag description").Int(1)
    
        app.Parse()
    
        // application logic ...
    }

In comparison to _flag_, _libgo-command_ provides additional configuration
options (e.g. name of flag arguments) and value parser/validator 
implementations (e.g. directories).

The intend of _libgo-command_ is not to replace _flag_, but to provide
a parser which categorizes the command-line in to more than just flags
and arguments.