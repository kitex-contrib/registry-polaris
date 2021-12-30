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

package polaris

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	perrors "github.com/pkg/errors"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"
)

const (
	defaultWeight           = 10
	PolarisDefaultNamespace = "default"
)

// Resolver is extension interface of Kitex Resolver.
type Resolver interface {
	discovery.Resolver

	doHeartbeat(ins *api.InstanceRegisterRequest)

	// todo add watch
}

// PolarisResolver is a resolver using polaris.
type PolarisResolver struct {
	namespace  string
	provider   api.ProviderAPI
	consumer   api.ConsumerAPI
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewPolarisResolver creates a polaris based resolver.
func NewPolarisResolver(endpoints []string) (Resolver, error) {
	sdkCtx, err := GetPolarisConfig(endpoints)
	if err != nil {
		return nil, perrors.WithMessage(err, "create polaris namingClient failed.")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	newInstance := &PolarisResolver{
		namespace:  PolarisDefaultNamespace,
		consumer:   api.NewConsumerAPIByContext(sdkCtx),
		provider:   api.NewProviderAPIByContext(sdkCtx),
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}

	return newInstance, nil
}

// Target implements the Resolver interface.
func (polaris *PolarisResolver) Target(ctx context.Context, target rpcinfo.EndpointInfo) (description string) {
	return target.ServiceName()
}

// Resolve implements the Resolver interface.
func (polaris *PolarisResolver) Resolve(ctx context.Context, desc string) (discovery.Result, error) {
	var (
		info instanceInfo
		eps  []discovery.Instance
	)

	getInstances := &api.GetInstancesRequest{}
	getInstances.Namespace = PolarisDefaultNamespace
	getInstances.Service = desc
	InstanceResp, err := polaris.consumer.GetInstances(getInstances)
	if nil != err {
		klog.Fatalf("fail to getOneInstance, err is %v", err)
	}
	instances := InstanceResp.GetInstances()
	if nil != instances {
		for _, instance := range instances {
			klog.Infof("instance getOneInstance is %s:%d", instance.GetHost(), instance.GetPort())
			weight := instance.GetWeight()
			if weight <= 0 {
				weight = defaultWeight
			}
			addr := instance.GetHost() + ":" + strconv.Itoa(int(instance.GetPort()))
			eps = append(eps, discovery.NewInstance(instance.GetProtocol(), addr, weight, info.Tags))
		}
	}

	if len(eps) == 0 {
		return discovery.Result{}, fmt.Errorf("no instance remains for %v", desc)
	}
	return discovery.Result{
		Cacheable: true,
		CacheKey:  desc,
		Instances: eps,
	}, nil
}

// Diff implements the Resolver interface.
func (polaris *PolarisResolver) Diff(cacheKey string, prev, next discovery.Result) (discovery.Change, bool) {
	return discovery.DefaultDiff(cacheKey, prev, next)
}

// Name implements the Resolver interface.
func (polaris *PolarisResolver) Name() string {
	return "Polaris"
}

// doHeartbeat Since polaris does not support automatic reporting of instance heartbeats, separate logic is needed to implement it.
func (polaris *PolarisResolver) doHeartbeat(ins *api.InstanceRegisterRequest) {
	ticker := time.NewTicker(5 * time.Second)

	heartbeat := &api.InstanceHeartbeatRequest{
		InstanceHeartbeatRequest: model.InstanceHeartbeatRequest{
			Service:   ins.Service,
			Namespace: ins.Namespace,
			Host:      ins.Host,
			Port:      ins.Port,
		},
	}

	for {
		select {
		case <-polaris.ctx.Done():
			return
		case <-ticker.C:
			polaris.provider.Heartbeat(heartbeat)
		}
	}
}
