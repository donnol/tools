package context

import (
	"context"

	"github.com/pkg/errors"
)

// 时间、地点、人物、事情
type (
	TimestampType  string
	RemoteAddrType string
	UserKeyType    string
	RequestKeyType string
)

const (
	// 时间
	TimestampKey TimestampType = "Timestamp"

	// 地点
	RemoteAddrKey RemoteAddrType = "RemoteAddr"

	// 用户
	UserKey UserKeyType = "UserID"

	// 请求
	RequestKey RequestKeyType = "RequestID"
)

// GetValue 从标准库ctx读取key对应value
func GetValue(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

// GetAllValue 获取ctx里的所有value
func GetAllValue(ctx context.Context) ([]interface{}, error) {
	r := make([]interface{}, 0, 4)

	r1, err := GetTimestampValue(ctx)
	if err != nil {
		return r, err
	}
	r = append(r, r1)

	r2, err := GetRemoteAddrValue(ctx)
	if err != nil {
		return r, err
	}
	r = append(r, r2)

	r3, err := GetUserValue(ctx)
	if err != nil {
		return r, err
	}
	r = append(r, r3)

	r4, err := GetRequestValue(ctx)
	if err != nil {
		return r, err
	}
	r = append(r, r4)

	return r, nil
}

func GetUserValue(ctx context.Context) (int, error) {
	v := GetValue(ctx, UserKey)
	vv, ok := v.(int)
	if !ok {
		return 0, errors.Errorf("get %s failed, got %v", UserKey, v)
	}
	return vv, nil
}

func GetRequestValue(ctx context.Context) (string, error) {
	v := GetValue(ctx, RequestKey)
	vv, ok := v.(string)
	if !ok {
		return "", errors.Errorf("get %s failed, got %v", RequestKey, v)
	}
	return vv, nil
}

func GetTimestampValue(ctx context.Context) (int64, error) {
	v := GetValue(ctx, TimestampKey)
	vv, ok := v.(int64)
	if !ok {
		return 0, errors.Errorf("get %s failed, got %v", TimestampKey, v)
	}
	return vv, nil
}

func GetRemoteAddrValue(ctx context.Context) (string, error) {
	v := GetValue(ctx, RemoteAddrKey)
	vv, ok := v.(string)
	if !ok {
		return "", errors.Errorf("get %s failed, got %v", RemoteAddrKey, v)
	}
	return vv, nil
}

func MustGetUserValue(ctx context.Context) int {
	v, err := GetUserValue(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func MustGetRequestValue(ctx context.Context) string {
	v, err := GetRequestValue(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func MustGetTimestampValue(ctx context.Context) int64 {
	v, err := GetTimestampValue(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func MustGetRemoteAddrValue(ctx context.Context) string {
	v, err := GetRemoteAddrValue(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
