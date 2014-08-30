package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var (
	// string notation lookup table for boolean values
	booleans map[string]bool
)

type Value interface {
	//Get() interface{}
	Set(string) error
}

// Boolean values can be infered by the presence of
// the flag on the command-line
type inferableValue interface {
	IsBoolFlag() bool
}

type boolValue struct {
	out *bool
}

type fileValue struct {
	out *string
	dir bool
}

type intValue struct {
	out *int
}

type voidValue struct {
	emulateBool bool
}

func (b *boolValue) Set(value string) error {
	if out, ok := booleans[value]; ok {
		*(b.out) = out

		return nil
	}

	return fmt.Errorf("'%s' is not a valid boolean value.", value)
}

func (b *boolValue) IsBoolFlag() bool {
	return true
}

func (f *fileValue) Set(value string) (err error) {
	if out, err := filepath.Abs(value); nil == err {
		if info, err := os.Stat(out); nil == err {
			// arrow-pattern alert!
			err = setFileValue(f.out, info, f.dir)
		}
	}

	return
}

func (i *intValue) Set(value string) (err error) {
	if out, err := strconv.Atoi(value); nil == err {
		*(i.out) = out
	}

	return
}

func (v *voidValue) Set(value string) error {
	return nil
}

func (v *voidValue) IsBoolFlag() bool {
	return v.emulateBool
}

func init() {
	booleans = make(map[string]bool)

	booleans[""] = true
	booleans["1"] = true
	booleans["y"] = true
	booleans["t"] = true
	booleans["T"] = true
	booleans["yes"] = true
	booleans["true"] = true
	booleans["TRUE"] = true
	booleans["True"] = true

	booleans["-1"] = false
	booleans["0"] = false
	booleans["n"] = false
	booleans["f"] = false
	booleans["F"] = false
	booleans["no"] = false
	booleans["false"] = false
	booleans["FALSE"] = false
	booleans["False"] = false
}

func setFileValue(out *string, meta os.FileInfo, mustDir bool) error {
	if meta.IsDir() {
		if false == mustDir {
			return fmt.Errorf("'%s' is a directory", meta.Name())
		}
	} else if mustDir {
		return fmt.Errorf("'%s' is not a directory", meta.Name())
	}

	*out = meta.Name()

	return nil
}
