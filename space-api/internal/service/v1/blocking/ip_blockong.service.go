package blocking

import (
	"fmt"
	"space-api/constants"
	"space-api/dto"
	"space-api/util"
	"space-api/util/arr"
	"space-api/util/performance"
	"space-domain/dao/biz"
	"space-domain/model"

	"golang.org/x/net/context"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type _ipBlockingService struct {
	cache performance.JsonCache
}

var DefaultIPBlockingService = &_ipBlockingService{
	cache: *performance.NewCache(constants.MB * 4),
}

func (s *_ipBlockingService) AddBlockingIP(req *dto.AddBlockingIPReq, ctx context.Context) (resp *dto.AddBlockingIPResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		ipTx := tx.BlockIPRecord
		_, e := ipTx.WithContext(ctx).
			Where(ipTx.IPAddr.Eq(req.IPAddr)).
			Take()

		// 不存在, 进行创建
		if e != nil {
			e = ipTx.WithContext(ctx).Create(&model.BlockIPRecord{
				IPAddr:      req.IPAddr,
				IPSource:    req.IPSource,
				UserAgent:   req.UserAgent,
				LastRequest: req.LastRequest,
			})

			if e != nil {
				return e
			}
		}

		return nil
	})
	if err != nil {
		err = util.CreateBizErr("添加阻止 IP 失败", err)
		return
	}
	resp = &dto.AddBlockingIPResp{}
	s.cache.Set(req.IPAddr, 1)

	return
}

func (s *_ipBlockingService) IpInBlockingList(ip string) bool {
	var tmp int64
	if err := s.cache.Get(ip, &tmp); err == nil {
		return true
	}

	return false
}

func (s *_ipBlockingService) GetBlockingPages(req *dto.GetBlockingPagesReq, ctx context.Context) (resp *dto.GetBlockingPagesResp, err error) {
	ipOp := biz.BlockIPRecord

	condList := []gen.Condition{}
	if req.WhereCondList != nil {
		for _, cond := range req.WhereCondList {
			parsed, e := cond.ParseCond(ipOp.TableName())
			if e != nil {
				err = util.CreateBizErr("参数错误: "+e.Error(), e)
				return
			}
			condList = append(condList, parsed)
		}
	}

	sortList := []field.Expr{}
	if req.SortCondList != nil {
		for _, s := range req.SortCondList {
			expr := s.ToOrderField(ipOp.TableName())
			sortList = append(sortList, expr)
		}
	}

	list, count, err := ipOp.WithContext(ctx).
		Where(condList...).
		Order(sortList...).
		FindByPage(req.Normalize())

	if err != nil {
		err = util.CreateBizErr("查询分页失败: "+err.Error(), err)
		return
	}

	resp = &dto.GetBlockingPagesResp{
		PageList: model.PageList[*model.BlockIPRecord]{
			List:  list,
			Page:  int64(*req.Page),
			Size:  int64(*req.Size),
			Total: count,
		},
	}

	return
}

func (s *_ipBlockingService) DeleteBlockingRecord(req *dto.DeleteBlockingRecordReq, ctx context.Context) (resp *dto.DeleteBlockingRecordResp, err error) {
	err = biz.Q.Transaction(func(tx *biz.Query) error {
		ipTx := tx.BlockIPRecord
		condList := []gen.Condition{}

		if req.WhereCondList != nil {
			for _, cond := range req.WhereCondList {
				parsed, e := cond.ParseCond(ipTx.TableName())
				if e != nil {
					return e
				}
				condList = append(condList, parsed)
			}
		}

		removes, e := ipTx.WithContext(ctx).
			Select(ipTx.ID).
			Where(condList...).
			Find()

		if e != nil {
			return e
		}

		_, e = ipTx.WithContext(ctx).
			Where(ipTx.ID.In(
				arr.MapSlice(
					removes,
					func(_ int, t *model.BlockIPRecord) int64 {
						return t.ID
					},
				)...,
			)).
			Delete()

		if e != nil {
			return e
		}

		// 启用异步任务区删除缓存的记录
		performance.DefaultTaskRunner.Go(func() {
			for _, r := range removes {
				s.cache.Delete(fmt.Sprintf("%d", r.ID))
			}
		})

		return nil
	})

	if err != nil {
		err = util.CreateBizErr("删除阻止记录失败", err)
		return
	}

	resp = &dto.DeleteBlockingRecordResp{}
	return
}

// SyncBlockingRecordInCache 同步数据库的阻止列表到缓存, 这个过程不是异步的
func (s *_ipBlockingService) SyncBlockingRecordInCache(ctx context.Context) (err error) {
	ipOp := biz.BlockIPRecord
	list, err := ipOp.WithContext(ctx).Find()
	if err != nil {
		err = util.CreateBizErr("获取列表失败", err)
	}
	s.cache.ClearAll()
	for _, ipRecord := range list {
		s.cache.Set(ipRecord.IPAddr, 1)
	}

	return
}
