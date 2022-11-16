package httpreq

import (
	"io"
	"net/http"
)

type HeaderAndReader struct {
	reader io.Reader

	header http.Header
}

var (
	_ io.Reader = (*HeaderAndReader)(nil)
)

func NewHeaderAndReader(
	r io.Reader,
	header http.Header,
) *HeaderAndReader {
	return &HeaderAndReader{
		reader: r,
		header: header,
	}
}

func (hr *HeaderAndReader) Read(p []byte) (n int, err error) {
	if hr.reader == nil {
		return 0, io.EOF
	}
	return hr.reader.Read(p)
}

func (hr *HeaderAndReader) Header() http.Header {
	return hr.header
}
