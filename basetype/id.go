package basetype

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Id 在json编码时会转为字符串，以防止js的数字溢出
type Id uint64

var (
	_ json.Marshaler   = (*Id)(nil)
	_ json.Unmarshaler = (*Id)(nil)
)

func (id Id) MarshalJSON() ([]byte, error) {
	s := strconv.FormatUint(uint64(id), 10)
	s = strconv.Quote(s)

	return []byte(s), nil
}

func (id *Id) UnmarshalJSON(data []byte) error {
	s := string(data)
	s, err := strconv.Unquote(s)
	if err != nil {
		return fmt.Errorf("[id] data unquote failed: %v of %s", err, data)
	}

	if s == "" {
		return nil
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fmt.Errorf("[id] data parse uint failed: %v of %s", err, data)
	}
	*id = Id(i)

	return nil
}
