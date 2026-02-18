//go:build !windows

package cache

import (
	"fmt"
	"os"
)

func checkNotRoot() error {
	if os.Geteuid() == 0 {
		return fmt.Errorf("cache operations should not run as root; use a regular user for goup run")
	}
	return nil
}
