package uv

import (
	"context"
	"space-api/dto"
	"space-api/util"
	"space-domain/dao/biz"
	"space-domain/model"
	"time"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type (
	IUVService interface {
		GetDailyUVCount(req *dto.GetDailyCountReq, ctx context.Context) (int64, error)
		QueryRangeUV(req *dto.QueryUvCountReq, ctx context.Context) (resp int64, err error)
		GetUVTrend(req *dto.GetUVTrendReq, ctx context.Context) ([]map[string]any, error)
		GetUvPages(req *dto.GetUvPagesReq, ctx context.Context) (resp *dto.GetUvPagesResp, err error)
		DeleteUVRecord(req *dto.DeleteUVReq, ctx context.Context) (resp *dto.DeleteUVResp, err error)
	}
	uvServiceImpl struct{}
)

var (
	_ IUVService = (*uvServiceImpl)(nil)

	DefaultUVService = &uvServiceImpl{}
)

// GetDailyUVCount 获取指定日期的独立访客数
func (m *uvServiceImpl) GetDailyUVCount(req *dto.GetDailyCountReq, ctx context.Context) (int64, error) {
	uvOp := biz.UVStatistic
	count, err := uvOp.
		WithContext(ctx).
		Where(uvOp.VisitDate.Eq(req.Date)).
		Count()

	if err != nil {
		err = util.CreateBizErr("获取访客数失败", err)
	}

	return count, err
}

// QueryRangeUV 获取指定日期范围内的独立访客数
func (m *uvServiceImpl) QueryRangeUV(req *dto.QueryUvCountReq, ctx context.Context) (resp int64, err error) {
	uvOp := biz.UVStatistic
	tableName := uvOp.TableName()
	condList := []gen.Condition{}
	if req.WhereCondList != nil {
		for _, cond := range req.WhereCondList {
			p, err := cond.ParseCond(tableName)
			if err != nil {
				return 0, util.CreateBizErr("查询参数错误: "+err.Error(), err)
			}
			condList = append(condList, p)

		}
	}
	resp, err = uvOp.WithContext(ctx).
		Where(condList...).
		Distinct(uvOp.VisitorHash).
		Count()
	if err != nil {
		return 0, util.CreateBizErr("统计出现错误", err)
	}

	return
}

// GetUVTrend 获取一段时间内的UV趋势
func (m *uvServiceImpl) GetUVTrend(req *dto.GetUVTrendReq, ctx context.Context) ([]map[string]any, error) {
	var results []struct {
		Date  string
		Count int64
	}
	uvOp := biz.UVStatistic

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -req.Days)

	err := uvOp.WithContext(ctx).
		Select(
			uvOp.VisitDate.As("date"),
			uvOp.VisitorHash.Distinct().Count().As("count"),
		).
		Where(
			uvOp.CreatedAt.Gte(startDate.UnixMilli()),
			uvOp.CreatedAt.Lte(endDate.UnixMilli()),
		).
		Group(uvOp.VisitDate).
		Order(uvOp.CreatedAt).
		Scan(&results)

	if err != nil {
		return nil, util.CreateBizErr("统计趋势失败", err)
	}

	// 转换为通用格式
	trend := make([]map[string]interface{}, len(results))
	for i, r := range results {
		trend[i] = map[string]interface{}{
			"date": r.Date,
			"uv":   r.Count,
		}
	}

	return trend, nil
}

func (m *uvServiceImpl) GetUvPages(req *dto.GetUvPagesReq, ctx context.Context) (resp *dto.GetUvPagesResp, err error) {
	uvOp := biz.UVStatistic
	tableName := uvOp.TableName()
	condList := []gen.Condition{}
	if req.WhereCondList != nil {
		for _, cond := range req.WhereCondList {
			p, err := cond.ParseCond(tableName)
			if err != nil {
				return nil, util.CreateBizErr("查询参数错误: "+err.Error(), err)
			}
			condList = append(condList, p)
		}
	}
	sortList := []field.Expr{}

	if req.SortList != nil {
		for _, o := range req.SortList {
			sortList = append(sortList, o.ToOrderField(tableName))
		}
	}

	list, count, err := uvOp.WithContext(ctx).
		Where(condList...).
		Order(sortList...).
		FindByPage(req.Normalize())
	if err != nil {
		err = util.CreateBizErr("查询错误", err)
		return
	}

	resp = &dto.GetUvPagesResp{
		PageList: model.PageList[*model.UVStatistic]{
			List:  list,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}
	return
}

func (m *uvServiceImpl) DeleteUVRecord(req *dto.DeleteUVReq, ctx context.Context) (resp *dto.DeleteUVResp, err error) {
	uvOp := biz.UVStatistic
	tableName := uvOp.TableName()
	condList := []gen.Condition{}
	if req.WhereCondList != nil {
		for _, cond := range req.WhereCondList {
			p, err := cond.ParseCond(tableName)
			if err != nil {
				return nil, util.CreateBizErr("查询参数错误: "+err.Error(), err)
			}
			condList = append(condList, p)
		}
	}
	_, err = uvOp.WithContext(ctx).
		Where(condList...).
		Delete()
	if err != nil {
		err = util.CreateBizErr("删除记录失败", err)
		return
	}

	resp = &dto.DeleteUVResp{}
	return
}
