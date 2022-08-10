// 通用组件配置

package cmd

import (
	"fmt"
)

// 最少需要参数数目
func MinimumNArgs(n int, args []string) error {
	if len(args) < n {
		return fmt.Errorf("requires at least %d arg(s), only received %d", n, len(args))
	}
	return nil
}

// 最多需要参数数目
func MaximumNArgs(n int, args []string) error {
	if len(args) > n {
		return fmt.Errorf("accepts at most %d arg(s), received %d", n, len(args))
	}
	return nil
}
