// Package timex provide location and some time helper function
package timex

import "time"

var (
	Location = time.FixedZone("CST", 8*3600) // 东八，Asia/Shanghai
)
