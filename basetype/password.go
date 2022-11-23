package basetype

import "golang.org/x/crypto/bcrypt"

type Password string

// Encrypt 使用bcrypt算法将明文密码哈希得到hash字符串
// bcrypt算法在对同一个密码哈希多次会得出不同结果，极大的保证了用户密码的安全
func (p Password) Encrypt() (pp string, err error) {
	r, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	return string(r), nil
}

// Compare 使用bcrypt算法判断密码是否与传入hash值匹配
func (p Password) Compare(hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}

// Password 实现String方法，打印时自动替换为*
func (p Password) String() string {
	return "*"
}
