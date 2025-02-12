package service

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"space-api/db"
	"space-api/dto"
	"space-api/dto/query"
	"space-api/util"
	"space-api/util/performance"
	"space-domain/dao/extra"
	"space-domain/model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huandu/xstrings"
	"gorm.io/gen/field"
)

type _logService struct {
}

var DefaultLogService = &_logService{}

func (*_logService) GetLogPages(req *dto.GetLogPagesReq, ctx *gin.Context) (resp *dto.GetLogPagesResp, err error) {
	logOp := extra.LogInfo
	condList, err := query.ParseCondList(logOp.TableName(), req.Conditions)

	if err != nil {
		err = util.CreateBizErr("参数错误: "+err.Error(), err)
		return
	}
	for _, col := range req.Conditions {
		if _, ok := logOp.GetFieldByName(xstrings.ToSnakeCase(col.Column)); !ok {
			err = util.CreateBizErr("非法的条件参数: "+col.Column, fmt.Errorf("illegal condition param: %s", col.Column))

			return
		}
	}

	op := logOp.WithContext(ctx).Where(condList...)
	// 存在排序字段
	if req.OrderColumns != nil {
		sorts := []field.Expr{}
		for _, col := range req.OrderColumns {
			if _, ok := logOp.GetFieldByName(xstrings.ToSnakeCase(col.Column)); !ok {
				err = util.CreateBizErr("非法的条件参数: "+col.Column, fmt.Errorf("illegal condition param: %s", col.Column))
				return
			}
			sorts = append(sorts, col.ToOrderField(logOp.TableName()))
		}

		op = op.Order(sorts...)
	}
	logPages, count, err := op.
		FindByPage(req.Normalize())
	if err != nil {
		err = util.CreateBizErr("查找分页数据失败", err)
		return
	}

	resp = &model.PageList[*model.LogInfo]{
		List:  logPages,
		Page:  int64(*req.Page),
		Size:  int64(*req.Size),
		Total: count,
	}

	return

}

func (*_logService) DeleteLogsByCondition(req *dto.DeleteLogReq, ctx *gin.Context) (resp *dto.DeleteLogResp, err error) {
	for _, col := range req.Conditions {
		if _, ok := extra.LogInfo.GetFieldByName(xstrings.ToSnakeCase(col.Column)); !ok {
			err = util.CreateBizErr("非法的条件参数: "+col.Column, fmt.Errorf("illegal condition param: %s", col.Column))

			return
		}
	}

	err = extra.Q.Transaction(func(tx *extra.Query) error {
		logOp := tx.LogInfo

		p, e := query.ParseCondList(logOp.TableName(), req.Conditions)
		if e != nil {
			err = util.CreateBizErr("参数错误: "+err.Error(), err)
			return e
		}
		_, e = logOp.WithContext(ctx).Where(p...).Delete()
		if e != nil {
			return e
		}

		return nil
	})
	if err != nil {
		err = util.CreateBizErr("删除日志错误", err)
		return
	}

	resp = &performance.Empty{}
	return
}

func (*_logService) DumLogsStream(req *dto.DumpLogReq, ctx *gin.Context) {
	rows, err := extra.LogInfo.WithContext(ctx).UnderlyingDB().Rows()
	if err != nil {
		ctx.Error(util.CreateBizErr("导出日志失败"+err.Error(), err))
		return
	}

	defer rows.Close()
	switch req.Format {
	case "json":
		outputJson(rows, ctx)
	case "markdown", "md":
		outputMarkdown(rows, ctx)
	case "csv":
		outputCSV(rows, ctx)
	default:
		ctx.Error(util.CreateBizErr("导出日志失败, 不支持的格式: "+req.Format, fmt.Errorf("un-support format: %s", req.Format)))
		return
	}
}

const timeFormat = "2006-01-02 15:04:05.000"

type dataTransfer = struct {
	ID            int64   `json:"id" xml:"id"`
	LogType       string  `json:"logType" xml:"logType"`
	Message       string  `json:"message" xml:"message"`
	Level         string  `json:"level" xml:"level"`
	CostTime      int64   `json:"costTime" xml:"costTime"`
	RequestMethod *string `json:"requestMethod" xml:"requestMethod"`
	RequestURI    *string `json:"requestURI" xml:"requestURI"`
	StackTrace    *string `json:"stackTrace" xml:"stackTrace"`
	IPAddr        *string `json:"ipAddr" xml:"ipAddr"`
	IPSource      *string `json:"ipSource" xml:"ipSource"`
	Useragent     *string `json:"useragent" xml:"useragent"`
	CreatedAt     string  `json:"createdAt" xml:"createdAt"`
}

func createConvertPtr(info *model.LogInfo) *dataTransfer {
	return &dataTransfer{
		ID:            info.ID,
		LogType:       info.LogType,
		Message:       info.Message,
		Level:         info.Level,
		CostTime:      info.CostTime,
		RequestMethod: info.RequestMethod,
		RequestURI:    info.RequestURI,
		StackTrace:    info.StackTrace,
		IPAddr:        info.IPAddr,
		IPSource:      info.IPSource,
		Useragent:     info.Useragent,
		CreatedAt:     time.UnixMilli(info.CreatedAt).Format(timeFormat),
	}
}

func outputJson(rows *sql.Rows, ctx *gin.Context) {
	ctx.Header("Content-Disposition", `attachment; filename="logs.json"`)
	ctx.Header("Content-Type", gin.MIMEJSON)
	db := db.GetExtraDB()
	ctx.Writer.WriteString("[\n")
	ls := list.New()
	for rows.Next() {
		info := &model.LogInfo{}
		db.ScanRows(rows, info)
		f, _ := json.Marshal(createConvertPtr(info))
		ls.PushBack(f)

		if ls.Len() >= 20 {
			for ls.Len() >= 5 {
				t := ls.Front()
				ls.Remove(t)
				ctx.Writer.WriteString("  ")
				ctx.Writer.Write(t.Value.([]byte))
				ctx.Writer.WriteString(",\n")
			}
		}
	}
	var l = ls.Len()
	for i := 0; i < l; i++ {
		// take
		t := ls.Front()
		ls.Remove(t)
		ctx.Writer.WriteString("  ")
		ctx.Writer.Write(t.Value.([]byte))
		ctx.Writer.WriteString(fmt.Sprintf("%s\n", util.TernaryExpr(l-1 == i, "", ",")))
	}
	ctx.Writer.WriteString("]\n")
	ctx.Writer.Flush()
}

func outputMarkdown(rows *sql.Rows, ctx *gin.Context) {
	ctx.Header("Content-Disposition", `attachment; filename="logs.md"`)
	ctx.Header("Content-Type", "text/markdown")
	db := db.GetExtraDB()
	ctx.Writer.WriteString("# 导出日志\n\n")

	valType := reflect.TypeOf(&dataTransfer{}).Elem()
	var headers []string
	for i := 0; i < valType.NumField(); i++ {
		headers = append(headers, valType.Field(i).Name)
	}

	var index = 0
	for rows.Next() {
		info := &model.LogInfo{}
		db.ScanRows(rows, info)
		val := reflect.ValueOf(createConvertPtr(info)).Elem()
		recordData := []string{}

		index++
		for i := 0; i < val.NumField(); i++ {
			recordData = append(
				recordData,
				fmt.Sprintf(
					"%s: %v%s",
					headers[i],
					val.Field(i),
					util.TernaryExpr(i != len(headers)-1, "; ", ""),
				),
			)
		}

		ctx.Writer.WriteString(fmt.Sprintf("%d. ", index) + strings.Join(recordData, ""))
		ctx.Writer.WriteString("\n")
	}
	ctx.Writer.Flush()
}

func outputCSV(rows *sql.Rows, ctx *gin.Context) {
	ctx.Header("Content-Disposition", `attachment; filename="logs.csv"`)
	ctx.Header("Content-Type", "text/csv")
	db := db.GetExtraDB()
	// Write CSV header (based on struct field names)
	val := reflect.TypeOf(&dataTransfer{}).Elem()
	var header []string
	for i := 0; i < val.NumField(); i++ {
		header = append(header, val.Field(i).Name)
	}
	ctx.Writer.WriteString(strings.Join(header, ", "))
	ctx.Writer.WriteString("\n")

	for rows.Next() {
		info := &model.LogInfo{}
		db.ScanRows(rows, info)

		val := reflect.ValueOf(createConvertPtr(info)).Elem()
		recordData := []string{}
		for i := 0; i < val.Type().NumField(); i++ {
			recordData = append(recordData, fmt.Sprintf("%v", val.Field(i)))
		}
		ctx.Writer.WriteString(strings.Join(recordData, ", "))
		ctx.Writer.WriteString("\n")
	}
	ctx.Writer.Flush()
}
