package rpeat

import (
	"strings"
)

// https://golangbot.com/custom-errors/

type Error int

const (
	OutOfCalRange Error = iota
)

func (e Error) String() string {
	err := [...]string{"Error: calendar out of range"}
	return err[e]
}

type ParseError struct {
	warn  bool
	errno []int
	err   []string
}

func (e ParseError) Error() string {
	return Stringify(e.err)
}

type specError struct {
	warn  bool
	errno int
	msg   string
}

type JobSpecErrors struct {
	errs []specError
	sep  string
}

func (e JobSpecErrors) Error() string {
	var msgs []string
	for _, err := range e.errs {
		msgs = append(msgs, err.msg)
	}
	return strings.Join(msgs, e.sep)
}
