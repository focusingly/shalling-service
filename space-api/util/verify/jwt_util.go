package verify

import (
	"fmt"
	"space-api/conf"
	"space-api/constants"
	"space-api/util/ptr"
	"space-domain/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var _jwtConf = conf.JwtConf{}

func init() {
	_jwtConf = *conf.ProjectConf.GetJwtConf()
}

type TokenParsedBizClaims struct {
	Iss string `json:"iss"`
	// 过期时间, 单位: 秒
	Exp int `json:"exp"`
	Iat int `json:"iat"`
	// 用户的 唯一ID
	Jti string `json:"jti"`
	// 用户类型标识
	UserType constants.UserType `json:"userType"`
}

func CreateJwtToken(loginSession *model.UserLoginSession) (token string, err error) {
	nowUnixSec := time.Now().Unix()
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "space.shalling.me",
		"exp":      int(loginSession.ExpiredAt / 1000),
		"nbf":      nowUnixSec,
		"iat":      nowUnixSec,
		"jti":      fmt.Sprintf("%d", loginSession.ID),
		"userType": loginSession.UserType,
	})

	return raw.SignedString(ptr.String2Bytes(_jwtConf.Salt))
}

func VerifyAndGetParsedBizClaims(tokenStr string) (parsedClaims *TokenParsedBizClaims, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return ptr.String2Bytes(_jwtConf.Salt), nil
	})
	if err != nil {
		return
	}
	if cl, ok := token.Claims.(jwt.MapClaims); !ok {
		err = fmt.Errorf("can't convert jwt claims")
		return
	} else {
		iss, err := cl.GetIssuer()
		if err != nil {
			return nil, err
		}
		exp, err := cl.GetExpirationTime()
		if err != nil {
			return nil, err
		}
		iat, err := cl.GetIssuedAt()
		if err != nil {
			return nil, err
		}
		parsedClaims = &TokenParsedBizClaims{
			Iss:      iss,
			Exp:      int(exp.Unix()),
			Iat:      int(iat.Unix()),
			Jti:      cl["jti"].(string),
			UserType: cl["userType"].(constants.UserType),
		}
	}

	return
}
