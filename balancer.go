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
	"strconv"
	"sync"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	polaris_go "github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
	"github.com/polarismesh/polaris-go/pkg/log"
	"github.com/polarismesh/polaris-go/pkg/model"
	"golang.org/x/sync/singleflight"
)

var polarisPickerPool sync.Pool

func init() {
	polarisPickerPool.New = newPolarisPicker
}

func newPolarisPicker() interface{} {
	return &polarisPicker{}
}

type polarisPicker struct {
	routerAPI           polaris_go.RouterAPI
	info                *polarisInfo
	routerInstancesResp *model.InstancesResponse
}

func (pp *polarisPicker) Next(ctx context.Context, request interface{}) (ins discovery.Instance) {
	svcInfo := model.ServiceInfo{
		Service:   pp.info.serviceName,
		Namespace: pp.info.namespace,
		// 目标服务的metadata，用于规则路由
		Metadata: pp.info.polarisOptions.DstMetadata,
	}

	routerRequest := &polaris_go.ProcessRoutersRequest{}
	routerRequest.DstInstances = model.NewDefaultServiceInstances(svcInfo, pp.info.polarisInstances)
	routerRequest.SourceService.Service = pp.info.polarisOptions.SrcService
	routerRequest.SourceService.Namespace = pp.info.polarisOptions.SrcNamespace
	routerRequest.SourceService.Metadata = pp.info.polarisOptions.SrcMetadata

	if pp.routerInstancesResp == nil {
		routerInstancesResp, err := pp.routerAPI.ProcessRouters(routerRequest)
		if nil != err {
			log.GetBaseLogger().Errorf("fail to do ProcessRouters err:%+v", err)
			return nil
		}
		if len(routerInstancesResp.GetInstances()) == 0 {
			return nil
		}
		pp.routerInstancesResp = routerInstancesResp
	}

	lbRequest := &polaris_go.ProcessLoadBalanceRequest{}
	lbRequest.DstInstances = pp.routerInstancesResp
	lbRequest.LbPolicy = config.DefaultLoadBalancerWR
	oneInstResp, err := pp.routerAPI.ProcessLoadBalance(lbRequest)
	if nil != err {
		return nil
	}

	targetInstance := oneInstResp.GetInstance()

	instanceKey := GetInstanceKey(targetInstance.GetNamespace(), targetInstance.GetService(), targetInstance.GetHost(), strconv.Itoa(int(targetInstance.GetPort())))

	// todo 没找到咋办
	ins = pp.info.polarisInstanceKeyKitexInstanceMap[instanceKey]

	return ins
}

func (pp *polarisPicker) Recycle() {
	pp.zero()
	polarisPickerPool.Put(pp)
}

func (pp *polarisPicker) zero() {
	pp.info = nil
	pp.routerAPI = nil
	pp.routerInstancesResp = nil
}

// polarisBalancer is a resolver using polaris.
type polarisBalancer struct {
	cachedPolarisInfo sync.Map
	sfg               singleflight.Group
	consumerAPI       api.ConsumerAPI
	routerAPI         polaris_go.RouterAPI

	polarisInstances []model.Instance
}

func NewPolarisBalancer(configFile ...string) (loadbalance.Loadbalancer, error) {
	sdkCtx, err := GetPolarisConfig(configFile...)
	if err != nil {
		return nil, err
	}

	pb := &polarisBalancer{
		consumerAPI: api.NewConsumerAPIByContext(sdkCtx),
		routerAPI:   polaris_go.NewRouterAPIByContext(sdkCtx),
	}

	return pb, nil
}

type polarisInfo struct {
	namespace                          string
	serviceName                        string
	kitexInstances                     []discovery.Instance
	polarisInstanceKeyKitexInstanceMap map[string]discovery.Instance
	polarisInstances                   []model.Instance
	polarisOptions                     Options
}

func (pb *polarisBalancer) GetPicker(e discovery.Result) loadbalance.Picker {
	var w *polarisInfo

	if e.Cacheable {
		cpi, ok := pb.cachedPolarisInfo.Load(e.CacheKey)
		if !ok {
			cpi, _, _ = pb.sfg.Do(e.CacheKey, func() (interface{}, error) {
				return pb.newPolarisInfo(e), nil
			})
			pb.cachedPolarisInfo.Store(e.CacheKey, cpi)
		}
		w = cpi.(*polarisInfo)
	} else {
		w = pb.newPolarisInfo(e)
	}

	picker := polarisPickerPool.Get().(*polarisPicker)
	picker.info = w
	picker.routerAPI = pb.routerAPI

	return picker
}

func (pb *polarisBalancer) Name() string {
	return "polaris"
}

// Rebalance implements the Rebalancer interface.
func (pb *polarisBalancer) Rebalance(change discovery.Change) {
	if !change.Result.Cacheable {
		return
	}
	pb.cachedPolarisInfo.Store(change.Result.CacheKey, pb.newPolarisInfo(change.Result))
}

// Delete implements the Rebalancer interface.
func (pb *polarisBalancer) Delete(change discovery.Change) {
	if !change.Result.Cacheable {
		return
	}
	pb.cachedPolarisInfo.Delete(change.Result.CacheKey)
}

func (pb *polarisBalancer) newPolarisInfo(e discovery.Result) *polarisInfo {
	pi := &polarisInfo{
		kitexInstances:                     make([]discovery.Instance, 0, len(e.Instances)),
		polarisInstances:                   make([]model.Instance, 0, len(e.Instances)),
		polarisInstanceKeyKitexInstanceMap: make(map[string]discovery.Instance, 0),
	}
	for _, kitexInst := range e.Instances {
		pkInst, ok := kitexInst.(*polarisKitexInstance)
		if !ok {
			continue
		}
		pi.polarisOptions = pkInst.polarisOptions
		pi.polarisInstances = append(pi.polarisInstances, pkInst.polarisInstance)
		pi.kitexInstances = append(pi.kitexInstances, kitexInst)

		instanceKey := GetInstanceKey(pkInst.polarisInstance.GetNamespace(), pkInst.polarisInstance.GetService(), pkInst.polarisInstance.GetHost(), strconv.Itoa(int(pkInst.polarisInstance.GetPort())))
		pi.polarisInstanceKeyKitexInstanceMap[instanceKey] = kitexInst
	}

	namespace, serviceName := SplitCachedKey(e.CacheKey)
	pi.namespace = namespace
	pi.serviceName = serviceName
	return pi
}
