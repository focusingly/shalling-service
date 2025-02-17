package service

import (
	"context"
	"space-api/conf"
	"space-api/util"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/accounts"
	"github.com/cloudflare/cloudflare-go/v3/option"
	"github.com/cloudflare/cloudflare-go/v3/shared"
)

type _clfService struct {
	*cloudflare.Client
	clfConf *conf.CloudflareConf
}

var DefaultCloudflareService *_clfService

func init() {
	clfConf := conf.ProjectConf.GetCloudflareConf()
	if clfConf == nil {
		return
	}
	client := cloudflare.NewClient(
		option.WithAPIKey(clfConf.ApiKey),  // defaults to os.LookupEnv("CLOUDFLARE_API_KEY")
		option.WithAPIEmail(clfConf.Email), // defaults to os.LookupEnv("CLOUDFLARE_EMAIL")
	)
	DefaultCloudflareService = &_clfService{
		clfConf: clfConf,
		Client:  client,
	}
}

// GetExistsCost 返回所有已经产生了费用的项目
func (s *_clfService) GetExistsCost(ctx context.Context) (subs []*shared.Subscription, err error) {
	list, err := s.Accounts.
		Subscriptions.
		Get(ctx, accounts.SubscriptionGetParams{
			AccountID: cloudflare.F(s.clfConf.AccountID),
		})
	if err != nil {
		err = util.CreateBizErr("查询订阅信息失败", err)
		return
	}

	// 返回所有已经产生费用的项目
	subs = []*shared.Subscription{}
	for _, sub := range *list {
		if sub.Price != 0 {
			subs = append(subs, &sub)
		}
	}

	return
}
