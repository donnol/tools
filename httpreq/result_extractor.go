package httpreq

import (
	"encoding/json"
	"encoding/xml"
)

type (
	ResultExtractor[R any] func(data []byte) (R, error)
)

func JSONExtractor[R any](data []byte) (R, error) {
	var r R
	if err := json.Unmarshal(data, &r); err != nil {
		return r, err
	}

	return r, nil
}

func XMLExtractor[R any](data []byte) (R, error) {
	var r R
	if err := xml.Unmarshal(data, &r); err != nil {
		return r, err
	}

	return r, nil
}
