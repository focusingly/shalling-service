package cmd

import (
	"space-api/db"
	"space-domain/dao/biz"
	"space-domain/dao/extra"
)

func setDataSource() {
	// 设置默认数据源
	biz.SetDefault(db.GetBizDB())
	extra.SetDefault(db.GetExtraDB())
}
