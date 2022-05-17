package debug

import (
	"fmt"
	"os"
)

func Printf(format string, args ...any) {
	if IsDebug() {
		fmt.Printf("| debug | "+format, args...)
	}
}

func IsDebug() bool {
	if v := os.Getenv("TOOLDEBUG"); v != "" {
		return true
	}
	return false
}
