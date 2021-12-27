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
	"net"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	perrors "github.com/pkg/errors"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"
)

var (
	protocolForKitex            string = "tcp"
	defaultHeartbeatIntervalSec        = 5
)

// Registry Extension - Registry
type Registry interface {
	registry.Registry

	doHeartbeat(ctx context.Context, ins *api.InstanceRegisterRequest)

	// todo add watch
}

type polarisRegistry struct {
	endpoints []string
	consumer  api.ConsumerAPI
	provider  api.ProviderAPI
}

// // NewPolarisRegistry creates a polaris based registry.
func NewPolarisRegistry(endpoints []string) (Registry, error) {

	sdkCtx, err := GetPolarisConfig(endpoints)
	if err != nil {
		return &polarisRegistry{}, err
	}
	pRegistry := &polarisRegistry{
		consumer: api.NewConsumerAPIByContext(sdkCtx),
		provider: api.NewProviderAPIByContext(sdkCtx),
	}

	return pRegistry, nil
}

// Register registers a server with given registry info.
func (svr *polarisRegistry) Register(info *registry.Info) error {
	if err := validateRegistryInfo(info); err != nil {
		return err
	}
	param := createRegisterParam(info)
	resp, err := svr.provider.Register(param)
	if err != nil {
		return err
	}
	if resp.Existed {
		klog.Warnf("instance already registered, namespace:%s, service:%s, port:%s",
			param.Namespace, param.Service, param.Host)
	}

	ctx, _ := context.WithCancel(context.Background())

	go svr.doHeartbeat(ctx, param)

	return nil
}

// Deregister deregisters a server with given registry info.
func (svr *polarisRegistry) Deregister(info *registry.Info) error {
	if info.ServiceName == "" {
		return fmt.Errorf("missing service name in Deregister")
	}

	request := createDeregisterParam(info)
	err := svr.provider.Deregister(request)
	if err != nil {
		return perrors.WithMessagef(err, "register(err:%+v)", err)
	}
	return nil

}

// IsAvailable always return true when use polaris
func (pr *polarisRegistry) IsAvailable() bool {
	return true
}

// doHeartbeat Since polaris does not support automatic reporting of instance heartbeats, separate logic is needed to implement it
func (svr *polarisRegistry) doHeartbeat(ctx context.Context, ins *api.InstanceRegisterRequest) {
	ticker := time.NewTicker(time.Duration(4) * time.Second)

	heartbeat := &api.InstanceHeartbeatRequest{
		InstanceHeartbeatRequest: model.InstanceHeartbeatRequest{
			Service:   ins.Service,
			Namespace: ins.Namespace,
			Host:      ins.Host,
			Port:      ins.Port,
			Timeout:   model.ToDurationPtr(60 * time.Second),
		},
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			svr.provider.Heartbeat(heartbeat)
		}
	}
}

func validateRegistryInfo(info *registry.Info) error {
	if info.ServiceName == "" {
		return fmt.Errorf("missing service name in Register")
	}
	if info.Addr == nil {
		return fmt.Errorf("missing addr in Register")
	}
	return nil
}

// createRegisterParam convert registry.Info to polaris instance register request
func createRegisterParam(info *registry.Info) *api.InstanceRegisterRequest {
	host, port, err := net.SplitHostPort(info.Addr.String())
	if err != nil {
		return nil
	}
	Instanceport, _ := strconv.Atoi(port)

	req := &api.InstanceRegisterRequest{
		InstanceRegisterRequest: model.InstanceRegisterRequest{
			Service:   info.ServiceName,
			Namespace: PolarisDefaultNamespace,
			Host:      host,
			Port:      Instanceport,
			Protocol:  &protocolForKitex,
			Timeout:   model.ToDurationPtr(10 * time.Second),
		},
	}

	req.SetTTL(defaultHeartbeatIntervalSec)

	return req
}

// createDeregisterParam convert registry.info to polaris instance deregister request
func createDeregisterParam(info *registry.Info) *api.InstanceDeRegisterRequest {
	host, port, err := net.SplitHostPort(info.Addr.String())
	if err != nil {
		return nil
	}
	Instanceport, _ := strconv.Atoi(port)

	return &api.InstanceDeRegisterRequest{
		InstanceDeRegisterRequest: model.InstanceDeRegisterRequest{
			Service:   info.ServiceName,
			Namespace: PolarisDefaultNamespace,
			Host:      host,
			Port:      Instanceport,
		},
	}
}