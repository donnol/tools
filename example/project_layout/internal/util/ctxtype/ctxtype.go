package ctxtype

import (
	"context"
	"fmt"
)

type (
	TimestampType  struct{}
	RemoteAddrType struct{}
	UserKeyType    struct{}
	RequestKeyType struct{}

	CheckDataPerm struct{}
	DataPermType  struct{}
	IsAdminType   struct{}
)

func GetUserID(ctx context.Context) (uint, error) {
	uid, ok := ctx.Value(UserKeyType{}).(uint)
	if !ok {
		return 0, fmt.Errorf("not login")
	}
	return uid, nil
}

func MustGetUserID(ctx context.Context) uint { //nolint
	uid, err := GetUserID(ctx)
	if err != nil {
		panic(err)
	}
	return uid
}

func GetTime(ctx context.Context) (int64, error) {
	data, ok := ctx.Value(TimestampType{}).(int64)
	if !ok {
		return 0, fmt.Errorf("not exist timestamp")
	}
	return data, nil
}

func MustGetTime(ctx context.Context) int64 { //nolint
	data, err := GetTime(ctx)
	if err != nil {
		panic(err)
	}
	return data
}

func NeedCheckPerm(ctx context.Context) bool {
	data := ctx.Value(CheckDataPerm{})
	v, ok := data.(bool)
	if ok && v {
		return true
	}
	return false
}

func GetDataPerm(ctx context.Context) interface{} {
	return ctx.Value(DataPermType{})
}

func IsAdmin(ctx context.Context) bool {
	data := ctx.Value(IsAdminType{})
	v, ok := data.(bool)
	if ok && v {
		return true
	}
	return false
}
