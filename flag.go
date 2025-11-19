package rpeat

import (
	"flag"
)

func argsPassed(flags *FlagNames) func(*flag.Flag) {
	return func(f *flag.Flag) {
		n := []string{f.Name}
		flags.Names = append(n, flags.Names...)
	}
}

type FlagNames struct {
	Names []string
}

func NewFlagNames() *FlagNames {
	return &FlagNames{}
}

func FlagArgs(flag *flag.FlagSet) []string {
	f := NewFlagNames()
	flag.Visit(argsPassed(f))
	return f.Names
}
