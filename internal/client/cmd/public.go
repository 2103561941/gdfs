// public tools

package cmd

import "fmt"

// check args length
func equalNumArgs(n int, args []string) error {
	if len(args) != n {
		return fmt.Errorf("requires %d arg(s), but received %d", n, len(args))
	}
	return nil
}
