// Package timex must import before all other packages, so `import _ "github.com/donnol/tools/timex"`` first, and put it a single block. And then,  put other imports to another block.
package timex

import "time"

var (
	Location = time.FixedZone("CST", 8*3600)
)

// init
func init() {
	time.Local = Location
}
