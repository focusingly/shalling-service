package service

import (
	"space-domain/dao/biz"

	"github.com/gin-gonic/gin"
)

func SelectTagListByPage(ctx *gin.Context) (val any, err error) {

	return
}

func DeleteTagByIdList(ctx *gin.Context) (val any, err error) {

	query := biz.Q

	query.Transaction(func(tx *biz.Query) error {
		return nil
	})
	return
}
