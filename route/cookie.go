package route

import (
	"net/http"
	"time"

	"github.com/donnol/tools/jwt"
)

type CookieOption struct {
	SessionKey string
	JwtToken   *jwt.Token
}

// MakeCookie 新建令牌
func MakeCookie(userID int, co CookieOption) (cookie http.Cookie, err error) {
	session, err := co.JwtToken.Sign(userID)
	if err != nil {
		return
	}

	days := 7
	var maxAge = 3600 * 24 * days

	cookie.Name = co.SessionKey
	cookie.Value = session
	cookie.MaxAge = maxAge
	cookie.Expires = time.Now().AddDate(0, 0, days)
	cookie.Path = "/"
	cookie.HttpOnly = true

	return
}
