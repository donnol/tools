package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"project_layout/internal/service"
	"project_layout/internal/util/ctxtype"
	"project_layout/model/db/user"
	"project_layout/model/request/common"

	"github.com/donnol/tools/jwt"
	"github.com/donnol/tools/route"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

type API struct {
	PingSrv service.PingSrv
	UserSrv service.UserSrv

	// TODO: add more field...
}

type (
	CookieOp string

	Param struct {
		method string

		// 参数
		body   []byte
		values url.Values

		CookieOp CookieOp
		UserID   uint
	}
)

const (
	CookieOpSet   CookieOp = "set"
	CookieOpUnset CookieOp = "unset"
)

type M = map[string]interface{}

// 参数相关
var (
	decoder = func() *schema.Decoder {
		dc := schema.NewDecoder()
		dc.IgnoreUnknownKeys(true)
		return dc
	}()
)

func NewParam(c *gin.Context) (*Param, error) {
	// 获取参数
	var err error
	var body []byte
	var values url.Values
	switch c.Request.Method {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		body, err = c.GetRawData()
	case http.MethodGet:
		fallthrough
	case http.MethodDelete:
		values = c.Request.URL.Query()
	}
	if err != nil {
		return nil, fmt.Errorf("new param failed: %w", err)
	}

	param := &Param{
		method: c.Request.Method,
		body:   body,
		values: values,
	}
	return param, nil
}

// Parse 解析
func (p *Param) Parse(ctx context.Context, v interface{}) error {
	var err error

	// 解析
	switch p.method {
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		err = json.Unmarshal(p.body, v)
	case http.MethodGet:
		fallthrough
	case http.MethodDelete:
		err = decoder.Decode(v, p.values)
	}
	if err != nil {
		return errors.WithStack(err)
	}

	// 检查参数
	if vv, ok := v.(Checker); ok {
		if err := vv.Check(ctx); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// Checker 检查接口
type Checker interface {
	Check(context.Context) error
}

// Filter 过滤器
type Filter interface {
	Filter() interface{}
}

type logicFunc = func(context.Context, *Param) (interface{}, error)

func jsonHandler(f logicFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取登录用户id
		userID := getUserIDByCookie(ctx)

		// ctx
		sctx, res := ctxFromGin(ctx, userID)

		// 业务
		param, res1 := handle(sctx, ctx, userID, f)
		res = res.Fill(res1)

		// 设置登录态
		setCookie(ctx, param)

		// 数据返回
		jsonResponse(ctx, param, http.StatusOK, res)
	}
}

func ctxFromGin(ctx *gin.Context, userID uint) (context.Context, common.Result) {
	var res common.Result

	// ctx
	now := time.Now().Unix()
	sctx := context.WithValue(ctx.Request.Context(), ctxtype.TimestampType{}, now)

	addr := ctx.Request.RemoteAddr
	sctx = context.WithValue(sctx, ctxtype.RemoteAddrType{}, addr)

	sctx = context.WithValue(sctx, ctxtype.UserKeyType{}, userID)

	traceID := uuid.New().String()
	sctx = context.WithValue(sctx, ctxtype.RequestKeyType{}, traceID)

	res.Timestamp = now
	res.TraceId = traceID

	return sctx, res
}

// func ctxWithPerms(ctx context.Context, userID uint) context.Context {
// 	// 获取用户数据权限
// 	perms, _ := srv.User.GetPerms(ctx, userID)
// 	ctx = context.WithValue(ctx, ctxtype.CheckDataPerm{}, true)
// 	ctx = context.WithValue(ctx, ctxtype.DataPermType{}, perms)
// 	return ctx
// }

func ctxWithIsAdmin(ctx context.Context) context.Context {
	_, err := checkAdmin(ctx)
	isAdmin := err == nil
	ctx = context.WithValue(ctx, ctxtype.IsAdminType{}, isAdmin)
	return ctx
}

func checkUser(ctx context.Context) (uint, error) {
	uid, err := ctxtype.GetUserID(ctx)
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func checkAdmin(ctx context.Context) (*user.Table, error) {
	uid, err := checkUser(ctx)
	if err != nil {
		return nil, err
	}
	_ = uid

	// ua, err := srv.User.GetWithAccess(ctx, uid)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func getUserIDByCookie(ctx *gin.Context) uint {
	var userID uint
	cookie, err := ctx.Cookie(sessionKey)
	if err == nil {
		verifyUserID, err := jwtToken.Verify(cookie)
		if err == nil {
			userID = uint(verifyUserID)
		}
	}
	return userID
}

func setCookie(ctx *gin.Context, param *Param) {
	const (
		setCookieHeader = "Set-Cookie"
	)
	switch param.CookieOp {
	case CookieOpSet:
		cookie, err := route.MakeCookie(int(param.UserID), route.CookieOption{
			SessionKey: sessionKey,
			JwtToken:   jwtToken,
		})
		if err != nil {
			panic(err)
		}
		ctx.Header(setCookieHeader, cookie.String())
	case CookieOpUnset:
		cookie := &http.Cookie{
			Name:     sessionKey,
			Value:    "",
			MaxAge:   0,
			Expires:  time.Now().AddDate(0, 0, -1),
			Path:     "/",
			HttpOnly: true,
		}
		ctx.Header(setCookieHeader, cookie.String())
	}
}

func jsonResponse(ctx *gin.Context, param *Param, code int, res common.Result) {
	ctx.Header("Content-Type", "application/json; charset=utf-8")

	ctx.JSON(code, res)
}

func handle(sctx context.Context, ctx *gin.Context, userID uint, f logicFunc) (*Param, common.Result) {
	var res common.Result

	// 参数
	param, err := NewParam(ctx)
	if err != nil {
		res.Code = 10001
		res.Msg = fmt.Sprintf("parse param failed: %v", err)
		return param, res
	}

	// 执行业务
	data, err := f(sctx, param)
	if err != nil {
		res.Code = 10002
		res.Msg = fmt.Sprintf("exec failed: %v", err)
		return param, res
	}

	if v, ok := data.(Filter); ok {
		data = v.Filter()
	}

	res.Code = 0
	res.Msg = ""
	res.Data = data

	return param, res
}

var (
	sessionKey string
	jwtToken   *jwt.Token
)

const (
	authHeaderKey  = "Authorization"
	authHeaderSign = "Signature"
)
