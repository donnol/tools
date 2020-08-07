package format

import "golang.org/x/tools/imports"

func Format(fileName, content string, formatOnly bool) (string, error) {
	options := &imports.Options{
		TabWidth:   8,
		TabIndent:  true,
		Comments:   true,
		Fragment:   true,
		FormatOnly: formatOnly,
	}
	formatContent, err := imports.Process(fileName, []byte(content), options)
	if err != nil {
		return "", err
	}

	return string(formatContent), nil
}
