package debug

import (
	"fmt"
	"os"
)

func Debug(format string, args ...any) {
	if v := os.Getenv("TOOLDEBUG"); v != "" {
		fmt.Printf("| debug | "+format, args...)
	}
}
