package util

import (
	"fmt"
	"log"
	"space-api/conf"
	"space-api/util/ptr"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtConf = conf.JwtConf{}

func init() {
	err := conf.GetProjectViper().UnmarshalKey("jwtConf", &jwtConf)
	if err != nil {
		log.Fatal(err)
	}
	if jwtConf.Expired.Value <= 0 {
		panic(fmt.Errorf("require a  positive value, but got: %d", jwtConf.Expired.Value))
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
		d = time.Second * time.Duration(jwtConf.Expired.Value)
	case "minute":
		d = time.Minute * time.Duration(jwtConf.Expired.Value)
	case "hour":
		d = time.Hour * time.Duration(jwtConf.Expired.Value)
	case "day":
		d = time.Hour * 24 * time.Duration(jwtConf.Expired.Value)
	default:
		panic(fmt.Errorf("un-support time unit: %s", jwtConf.Expired.Unit))
	}

	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "shalling.me",
		"sub":  "auth-token",
		"aud":  []string{"visit", "comment"},
		"exp":  now.Add(d).Unix(),
		"nbf":  now.Unix(),
		"iat":  now.Unix(),
		"jti":  id,
		"uuid": uuid.NewString(),
	})

	return raw.SignedString(ptr.String2Bytes(jwtConf.Salt))
}

func VerifyJwtToken(tokenStr string) (id string, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return ptr.String2Bytes(jwtConf.Salt), nil
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

func VerifyAndGetClaims(tokenStr string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return ptr.String2Bytes(jwtConf.Salt), nil
	})

	if err != nil {
		return
	}
	if cl, ok := token.Claims.(jwt.MapClaims); !ok {
		err = fmt.Errorf("can't convert jwt claims")

		return
	} else {
		claims = cl
	}
	return
}
