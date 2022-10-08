package httpreq

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Do[R any](
	client *http.Client,
	method string,
	link string,
	body io.Reader,
	codeChecker CodeChecker,
	extractResult ResultExtractor[R],
) (R, error) {
	var r R

	if method == "" || link == "" {
		return r, fmt.Errorf("bad param")
	}

	req, err := http.NewRequest(method, link, body)
	if err != nil {
		return r, err
	}

	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	err = codeChecker(resp.StatusCode)
	if err != nil {
		return r, fmt.Errorf("check code failed: %v, data: %s", err, data)
	}

	r, err = extractResult(data)
	if err != nil {
		return r, fmt.Errorf("extract result failed: %v, data: %s", err, data)
	}

	return r, nil
}
