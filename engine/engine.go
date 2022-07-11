package engine

import (
	"fmt"
	"os"
)

// Run takes a set of script files, parses them and runs resulting scripts.
// Each script is run right after it is parsed.
// A single evaluation context is used for all script files, i.e. definitions and results of the former
// script files may be shadowed (overridden) by definitions of the same name (functions and variables)
// in the latter script files.
func Run(scriptFiles ...string) error {
	ctx := newContext(nil)

	for _, f := range scriptFiles {
		err := func() error {
			file, err := os.Open(f)
			if err != nil {
				return err
			}

			defer func() {
				if err := file.Close(); err != nil {
					fmt.Printf("cannot close file %q: %v", f, err)
				}
			}()

			s, err := parseScript(file)
			if err != nil {
				return err
			}

			if err = s.eval(ctx); err != nil {
				return err
			}

			return nil
		}()

		if err != nil {
			return fmt.Errorf(
				"aborting execution: failed to run script %q: %v", f, err,
			)
		}
	}

	return nil
}
