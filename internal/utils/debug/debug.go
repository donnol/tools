package debug

import (
	"fmt"
	"os"
)

func Debug(format string, args ...interface{}) {
	if v := os.Getenv("TOOLDEBUG"); v != "" {
		fmt.Printf("| debug | "+format, args...)
	}
}
