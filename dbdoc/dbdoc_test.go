package dbdoc

import (
	"os"
)

func OpenFile(file string) (*os.File, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return f, err
	}

	_, err = f.WriteString("# 数据库表定义\n\n")
	if err != nil {
		return f, err
	}

	return f, nil
}
