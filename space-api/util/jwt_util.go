package util

import (
	"fmt"
	"log"
	"space-api/conf"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtConf = conf.JwtConf{}

func init() {
	err := conf.GetProjectViper().UnmarshalKey("jwtConf", &jwtConf)
	if err != nil {
		log.Fatal(err)
	}
	if jwtConf.Expired.Setup <= 0 {
		panic(fmt.Errorf("require a  positive value, but got: %d", jwtConf.Expired.Setup))
	}
	switch jwtConf.Expired.Unit {
	case "second",
		"minute",
		"hour",
		"day":
		return
	default:
		panic(fmt.Errorf("un-support time unit: %s", jwtConf.Expired.Unit))
	}
}

func CreateJwtToken(id string) (token string, err error) {
	now := time.Now()
	var d time.Duration

	switch jwtConf.Expired.Unit {
	case "second":
		d = time.Second * time.Duration(jwtConf.Expired.Setup)
	case "minute":
		d = time.Minute * time.Duration(jwtConf.Expired.Setup)
	case "hour":
		d = time.Hour * time.Duration(jwtConf.Expired.Setup)
	case "day":
		d = time.Hour * 24 * time.Duration(jwtConf.Expired.Setup)
	default:
		panic(fmt.Errorf("un-support time unit: %s", jwtConf.Expired.Unit))
	}

	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:   "shalling.me",
		Subject:  "The",
		Audience: []string{"user"},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(d),
		},
		NotBefore: &jwt.NumericDate{
			Time: now,
		},
		IssuedAt: &jwt.NumericDate{
			Time: now,
		},
		ID: id,
	})
	return raw.SignedString(jwtConf.Salt)
}

func VerifyJwtToken(tokenStr string) (id string, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtConf.Salt, nil
	})
	if err != nil {
		return
	}

	if cl, ok := token.Claims.(jwt.MapClaims); ok {
		if id, ok2 := cl["jti"].(string); ok2 {
			return id, nil
		}

	} else {
		err = fmt.Errorf("can't extract id value")
	}

	return
}
