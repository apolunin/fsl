package engine

import (
	"errors"
	"fmt"
	"strings"
)

// command represents an FSL command with a list of arguments
type command struct {
	cmd       string
	arguments []argument
}

// eval evaluates a command.
// The idea is to create a new context (a stack frame) with current context as its parent,
// evaluate command args in the parent context and push them to a new context. New context
// with pushed args is then used to evaluate a command.
func (c *command) eval(ctx *context) error {
	f, ok := ctx.getFunc(c.cmd, true)
	if !ok {
		return fmt.Errorf(
			"failed to evaluate command: function %q is not defined", c.cmd,
		)
	}

	newCtx := newContext(ctx)

	for _, arg := range c.arguments {
		newCtx.variables[arg.name()] = arg.value(ctx)
	}

	return f.eval(newCtx)
}

// parseCommand creates a command from unstructured JSON
func parseCommand(data map[string]any) (*command, error) {
	c, ok := data["cmd"]
	if !ok {
		return nil, errors.New(
			"cannot parse command: 'cmd' attribute is missing",
		)
	}

	var args []argument

	for k, v := range data {
		if k == "cmd" {
			continue
		}

		arg, err := parseArg(k, fmt.Sprintf("%v", v))
		if err != nil {
			return nil, fmt.Errorf(
				"cannot parse arg %q: %v", k, err,
			)
		}

		args = append(args, arg)
	}

	cmd := fmt.Sprintf("%v", c)
	if strings.HasPrefix(cmd, "#") {
		cmd = cmd[1:]
	}

	return &command{
		cmd:       cmd,
		arguments: args,
	}, nil
}
