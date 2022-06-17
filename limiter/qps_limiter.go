/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package limiter

import (
	"context"
	"time"

	polaris "github.com/kitex-contrib/registry-polaris"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/log"
)

// qpsLimiter is a gRPC interceptor that implements rate limiting.
type qpsLimiter struct {
	namespace string
	svcName   string
	limitAPI  api.LimitAPI
}

// NewQPSLimiter creates a new qpsLimiter.
func NewQPSLimiter(configFile ...string) (*qpsLimiter, error) {
	sdkCtx, err := polaris.GetPolarisConfig(configFile...)
	if err != nil {
		return nil, err
	}

	return &qpsLimiter{limitAPI: api.NewLimitAPIByContext(sdkCtx)}, nil
}

// WithNamespace sets the namespace of the service.
func (p *qpsLimiter) WithNamespace(namespace string) *qpsLimiter {
	p.namespace = namespace
	return p
}

// WithServiceName sets the service name.
func (p *qpsLimiter) WithServiceName(svcName string) *qpsLimiter {
	p.svcName = svcName
	return p
}

func (p *qpsLimiter) Acquire(ctx context.Context) bool {
	quotaReq := api.NewQuotaRequest()
	quotaReq.SetNamespace(p.namespace)
	quotaReq.SetService(p.svcName)
	future, err := p.limitAPI.GetQuota(quotaReq)
	if nil != err {
		log.GetBaseLogger().Errorf("fail to do ratelimit, err is %v", err)
		return false
	}
	rsp := future.Get()
	if rsp.Code != api.QuotaResultOk {
		return false
	}
	return true
}

func (p *qpsLimiter) Status(ctx context.Context) (max, current int, interval time.Duration) {
	return
}
