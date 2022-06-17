package limiter

import (
	"context"

	"github.com/cloudwego/kitex/pkg/limiter"
	polaris "github.com/kitex-contrib/registry-polaris"
	"github.com/polarismesh/polaris-go/api"
)

// polarisConcurrencyLimiter implements ConcurrencyLimiter.
type polarisConcurrencyLimiter struct {
	lim int32
	now int32
	tmp int32

	namespace   string
	serviceName string
	limitAPI    api.LimitAPI
}

// NewPolarisConcurrencyLimiter returns a new ConcurrencyLimiter with the given limit.
func NewPolarisConcurrencyLimiter(namespace string, serviceName string, configFile ...string) (limiter.ConcurrencyLimiter, error) {
	sdkCtx, err := polaris.GetPolarisConfig(configFile...)
	if err != nil {
		return nil, err
	}
	pcl := &polarisConcurrencyLimiter{
		lim:         1,
		now:         0,
		tmp:         0,
		namespace:   namespace,
		serviceName: serviceName,
		limitAPI:    api.NewLimitAPIByContext(sdkCtx),
	}
	return pcl, nil
}

func (pcl *polarisConcurrencyLimiter) Acquire(ctx context.Context) bool {
	quotaReq := api.NewQuotaRequest()
	quotaReq.SetNamespace(pcl.namespace)
	quotaReq.SetService(pcl.serviceName)
	future, err := pcl.limitAPI.GetQuota(quotaReq)
	if nil != err {
		// todo，err的处理
		return false
	}
	rsp := future.Get()
	if rsp.Code == api.QuotaResultLimited {
		return false
	}
	return true
}

func (pcl *polarisConcurrencyLimiter) Release(ctx context.Context) {
	quotaReq := api.NewQuotaRequest()
	quotaReq.SetNamespace(pcl.namespace)
	quotaReq.SetService(pcl.serviceName)
	future, err := pcl.limitAPI.GetQuota(quotaReq)
	if nil != err {
		return
	}
	future.Release()
}

func (pcl *polarisConcurrencyLimiter) Status(ctx context.Context) (limit, occupied int) {
	return
}
