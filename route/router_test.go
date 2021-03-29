package route

import (
	"context"
	"testing"
)

func getUser(ctx context.Context, p Param) (r Result, err error) {

	return
}

func getUserCurrent(ctx context.Context, p Param) (r Result, err error) {

	return
}

func TestGetMethodPathFromFunc(t *testing.T) {
	method, path := getMethodPathFromFunc(getUser)
	t.Log(method, path)

	method, path = getMethodPathFromFunc(getUserCurrent)
	t.Log(method, path)
}
