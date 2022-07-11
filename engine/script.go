package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

// script represents a parsed FSL script
type script struct {
	functions map[string]function
	variables map[string]string
}

// eval evaluates a script.
// It copies all script functions and variables to the provided evaluation context and
// evaluates 'init' function in this context.
func (s *script) eval(ctx *context) error {
	for key, value := range s.functions {
		ctx.functions[key] = value
	}

	for key, value := range s.variables {
		ctx.variables[key] = value
	}

	return s.functions[funcInit].eval(ctx)
}

// parseScript parses a script, reading its contents from io.Reader.
func parseScript(r io.Reader) (*script, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read data: %v", err)
	}

	var data map[string]any
	if err = json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("cannot parse JSON: %v", err)
	}

	var (
		functions = make(map[string]function)
		variables = make(map[string]string)
	)

	addBuiltIns(functions)

	for key, val := range data {
		switch value := val.(type) {
		case []any:
			f, err := parseFunction(key, value)
			if err != nil {
				return nil, fmt.Errorf("cannot parse script: %v", err)
			}
			functions[key] = f
		case float32, float64, int8, int16, int32, int64, int:
			v := fmt.Sprintf("%v", value)
			variables[key] = v
		case string:
			// As per requirements, all variables should be floating-point numbers.
			// Let's try to perform this conversion even if variable is of type 'string' in JSON.
			_, err = strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf(
					"variable %q = %q is not a floating-point number", key, value,
				)
			}
			variables[key] = value
		}
	}

	if _, ok := functions[funcInit]; !ok {
		return nil, fmt.Errorf("%q function not found", funcInit)
	}

	return &script{
		functions: functions,
		variables: variables,
	}, nil
}
