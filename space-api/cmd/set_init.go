package cmd

import (
	"context"
	"log"
	"space-api/db"
	"space-api/util/encrypt"
	"space-domain/dao/biz"
	"space-domain/model"
)

func prepareStartup() {
	q := biz.Q.ReplaceDB(db.GetBizDB())
	_, err := q.LocalUser.WithContext(context.TODO()).Where(q.LocalUser.IsAdmin.Gt(0)).Take()
	if err != nil {
		hashedPass, e := encrypt.EncryptPasswordByBcrypt("12345678")
		if e != nil {
			log.Fatal("init default password failed", e)
		}
		e = q.LocalUser.WithContext(context.TODO()).Create(&model.LocalUser{
			Username:    "shalling-admin233",
			DisplayName: "Shalling's Space",
			Password:    hashedPass,
			IsAdmin:     1,
		})
		if e != nil {
			log.Fatal("create default user failed", e)
		}
	}
}
