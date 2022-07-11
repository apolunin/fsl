package engine

import (
	"strings"
)

var _ argument = (*valueArg)(nil)
var _ argument = (*paramArg)(nil)
var _ argument = (*refArg)(nil)

const undefinedValue = "undefined"

// argument represents any kind of arg which can be passed to command.
type argument interface {
	name() string
	value(ctx *context) string
}

// parseArg parses an argument value and instantiates appropriate implementation.
func parseArg(name, value string) (argument, error) {
	switch {
	case strings.HasPrefix(value, "$"):
		return &paramArg{label: name, ref: value[1:]}, nil
	case strings.HasPrefix(value, "#"):
		return &refArg{label: name, ref: value[1:]}, nil
	}

	return &valueArg{label: name, val: value}, nil
}

// valueArg represents an argument without indirection (i.e. arg value has no '#' or '$' prefixes)
type valueArg struct {
	label string
	val   string
}

func (arg *valueArg) name() string {
	return arg.label
}

func (arg *valueArg) value(_ *context) string {
	return arg.val
}

// paramArg represents an argument to a function call (i.e. arg value has '$' prefix).
// It should be in the current context, i.e. paramArg is not searched in parent contexts.
type paramArg struct {
	label string
	ref   string
}

func (arg *paramArg) name() string {
	return arg.label
}

func (arg *paramArg) value(ctx *context) string {
	if val, ok := ctx.getVar(arg.ref, false); ok {
		return val
	}

	return undefinedValue
}

// refArg represents a reference to a variable (i.e. arg value has '#' prefix).
// refArg is searched in the current context. If it is not found, it is searched in the parent context.
// The process continues until root context is reached or value is found.
type refArg struct {
	label string
	ref   string
}

func (arg *refArg) name() string {
	return arg.label
}

func (arg *refArg) value(ctx *context) string {
	if val, ok := ctx.getVar(arg.ref, true); ok {
		return val
	}

	return undefinedValue
}
