package engine

// context represents an evaluation context for FSL script.
type context struct {
	parent *context

	functions map[string]function
	variables map[string]string
}

// getVar searches for a variable by name in context.
// If 'searchParent' is false, only current context is considered,
// otherwise parent context(s) are searched as well.
// If a variable is found, the function returns its value and true to communicate successful search,
// otherwise empty string and false is returned to communicate that variable is not found.
func (ctx *context) getVar(name string, searchParent bool) (string, bool) {
	if !searchParent {
		val, ok := ctx.variables[name]
		return val, ok
	}

	for curr := ctx; curr != nil; curr = curr.parent {
		if val, ok := curr.variables[name]; ok {
			return val, ok
		}
	}

	return "", false
}

// setVar sets a variable in a context.
// It starts from current context and tries to search for a variable.
// If a variable is found in a current context then its value is updated,
// otherwise the same process repeats in a parent context.
// If a variable doesn't exist in any context, a brand new one is created in a root context.
func (ctx *context) setVar(name, value string) {
	var curr *context

	for curr = ctx; curr.parent != nil; curr = curr.parent {
		if _, ok := curr.variables[name]; ok {
			curr.variables[name] = value
			return
		}
	}

	curr.variables[name] = value
}

// updateVar updates a variable in a context.
// It starts from current context and tries to search for a variable.
// If a variable is found in a current context then its value is updated,
// otherwise the same process repeats in a parent context.
// If a variable doesn't exist in any context, nothing happens and the function returns false,
// otherwise a variable's value is updated and the function returns true.
func (ctx *context) updateVar(name, value string) bool {
	for curr := ctx; curr != nil; curr = curr.parent {
		if _, ok := curr.variables[name]; ok {
			curr.variables[name] = value
			return true
		}
	}

	return false
}

// deleteVar deletes a variable from a context.
// It works the same way as updateVar, but deletes a variable instead of updating it.
func (ctx *context) deleteVar(name string) bool {
	for curr := ctx; curr != nil; curr = curr.parent {
		if _, ok := curr.variables[name]; ok {
			delete(curr.variables, name)
			return true
		}
	}

	return false
}

// getFunc searches for a function by name in a context.
// It works the same way as getVar, but works with functions instead of variables.
func (ctx *context) getFunc(name string, searchParent bool) (function, bool) {
	if !searchParent {
		f, ok := ctx.functions[name]
		return f, ok
	}

	for curr := ctx; curr != nil; curr = curr.parent {
		if f, ok := curr.functions[name]; ok {
			return f, ok
		}
	}

	return nil, false
}

// newContext creates a new context and sets its parent context.
func newContext(parent *context) *context {
	return &context{
		parent:    parent,
		functions: make(map[string]function),
		variables: make(map[string]string),
	}
}
