package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CustomClaims struct {
	AccountNum string `json:"account_num"`
	Password   string `json:"password"`
	jwt.RegisteredClaims
}

const TokenExpireDuration = time.Hour * 1 //1小时过期

// CustomSecret 用于加密的字符串
var CustomSecret = []byte("canyang")

// GenToken 生成JWT
func GenToken(password string, accountNum string) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		accountNum, // 自定义字段
		password,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "my-project", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
