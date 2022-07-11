package engine

import (
	"fmt"
	"strconv"
)

const (
	funcCreate = "create"
	funcDelete = "delete"
	funcUpdate = "update"
	funcPrint  = "print"
	funcInit   = "init"
	funcAdd    = "add"
	funcSub    = "sub"
	funcMul    = "mul"
	funcDiv    = "div"
)

const (
	attrID  = "id"
	attrOp1 = "operand1"
	attrOp2 = "operand2"
	attrVal = "value"
)

// function represents FSL function (custom or built-in) which can be evaluated.
type function interface {
	eval(ctx *context) error
}

// customFunc represents a user-defined FSL function (a list of commands).
type customFunc struct {
	name     string
	commands []*command
}

func (f *customFunc) eval(ctx *context) error {
	for _, command := range f.commands {
		if err := command.eval(ctx); err != nil {
			return fmt.Errorf(
				"failed to execute function %q: %v", f.name, err,
			)
		}
	}

	return nil
}

// parseFunction creates a user-defined FSL function from slice of commands.
func parseFunction(name string, data []any) (function, error) {
	var commands []*command
	for _, cmd := range data {
		switch val := cmd.(type) {
		case map[string]any:
			command, err := parseCommand(val)
			if err != nil {
				return nil, fmt.Errorf(
					"cannot parse function %q: %v", name, err,
				)
			}

			commands = append(commands, command)
		default:
			return nil, fmt.Errorf(
				"cannot parse function %q: invalid command: %v", name, val,
			)
		}
	}

	return &customFunc{
		name:     name,
		commands: commands,
	}, nil
}

// addBuiltIns adds built-in functions to the functions map.
func addBuiltIns(functions map[string]function) {
	functions[funcCreate] = &createFunc{name: funcCreate}
	functions[funcDelete] = &deleteFunc{name: funcDelete}
	functions[funcUpdate] = &updateFunc{name: funcUpdate}
	functions[funcPrint] = &printFunc{name: funcPrint}

	functions["add"] = &binOp{
		name: funcAdd,
		op:   func(a, b float64) float64 { return a + b },
	}
	functions["sub"] = &binOp{
		name: funcSub,
		op:   func(a, b float64) float64 { return a - b },
	}
	functions["mul"] = &binOp{
		name: funcMul,
		op:   func(a, b float64) float64 { return a * b },
	}
	functions["div"] = &binOp{
		name: funcDiv,
		op:   func(a, b float64) float64 { return a / b },
	}
}

// createFunc represents a built-in function 'create'.
type createFunc struct{ name string }

// deleteFunc represents a built-in function 'delete'.
type deleteFunc struct{ name string }

// updateFunc represents a built-in function 'update'.
type updateFunc struct{ name string }

// printFunc represents a built-in function 'print'.
type printFunc struct{ name string }

// binOp represents a binary operator. It is used to implement built-in
// functions 'add', 'sub', 'mul' and 'div'.
type binOp struct {
	name string
	op   func(float64, float64) float64
}

func (f *binOp) eval(ctx *context) error {
	id, ok := ctx.getVar(attrID, false)
	if !ok {
		return fmt.Errorf(
			"cannot evaluate %q function: %q is missing", f.name, attrID,
		)
	}

	op1, ok := ctx.getVar(attrOp1, false)
	if !ok {
		return fmt.Errorf(
			"cannot evaluate %q function: %q is missing", f.name, attrOp1,
		)
	}

	op2, ok := ctx.getVar(attrOp2, false)
	if !ok {
		return fmt.Errorf(
			"cannot evaluate %q function: %q is missing", f.name, attrOp2,
		)
	}

	left, err := strconv.ParseFloat(op1, 64)
	if err != nil {
		return fmt.Errorf(
			"cannot convert %q = %q to float: %v", attrOp1, op1, err,
		)
	}

	right, err := strconv.ParseFloat(op2, 64)
	if err != nil {
		return fmt.Errorf(
			"cannot convert %q = %q to float: %v", attrOp2, op2, err,
		)
	}

	res := f.op(left, right)
	ctx.setVar(id, fmt.Sprintf("%f", res))

	return nil
}

func (f *createFunc) eval(ctx *context) error {
	id, ok := ctx.getVar(attrID, false)
	if !ok {
		return fmt.Errorf(
			"failed to evaluate function %q: missing %q argument", f.name, attrID,
		)
	}

	val, ok := ctx.getVar(attrVal, false)
	if !ok {
		return fmt.Errorf(
			"failed to evaluate function %q: missing %q argument", f.name, attrVal,
		)
	}

	ctx.setVar(id, val)
	return nil
}

func (f *deleteFunc) eval(ctx *context) error {
	id, ok := ctx.getVar(attrID, false)
	if !ok {
		return fmt.Errorf(
			"failed to evaluate function %q: missing %q argument", f.name, attrID,
		)
	}

	if !ctx.deleteVar(id) {
		return fmt.Errorf(
			"failed to evaluate function %q: delete failed, variable %q is undefined", f.name, id,
		)
	}

	return nil
}

func (f *updateFunc) eval(ctx *context) error {
	id, ok := ctx.getVar(attrID, false)
	if !ok {
		return fmt.Errorf(
			"failed to evaluate function %q: missing %q argument", f.name, attrID,
		)
	}

	val, ok := ctx.getVar(attrVal, false)
	if !ok {
		return fmt.Errorf(
			"failed to evaluate function %q: missing %q argument", f.name, attrVal,
		)
	}

	if !ctx.updateVar(id, val) {
		return fmt.Errorf(
			"failed to evaluate function %q: update failed, variable %q is undefined", f.name, id,
		)
	}

	return nil
}

func (f *printFunc) eval(ctx *context) error {
	val, ok := ctx.getVar(attrVal, false)
	if !ok {
		return fmt.Errorf(
			"failed to evaluate function %q: missing %q argument", f.name, attrVal,
		)
	}

	if num, err := strconv.ParseFloat(val, 64); err == nil {
		fmt.Printf("%.4f\n", num)
	} else {
		fmt.Println(val)
	}

	return nil
}
