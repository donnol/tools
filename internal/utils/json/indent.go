package json

import (
	"bytes"
	"encoding/json"
	"os"
)

func IndentToStdout(r any) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := json.Indent(buf, data, "", "    "); err != nil {
		return err
	}
	if _, err = buf.WriteTo(os.Stdout); err != nil {
		return err
	}

	return nil
}
