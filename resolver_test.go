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
	"testing"
	"time"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/stretchr/testify/require"
)

const (
	serviceName = "registry-polaris"
)

func TestPolarisResolver(t *testing.T) {
	rg, err := NewPolarisRegistry([]string{"127.0.0.1:8091"})
	require.Nil(t, err)
	rs, err := NewPolarisResolver([]string{"127.0.0.1:8091"})
	require.Nil(t, err)

	// test register service
	info0 := &registry.Info{
		ServiceName: serviceName,
		Addr:        utils.NewNetAddr("tcp", "127.0.0.1:6666"),
		Weight:      100,
		Tags:        nil,
	}
	err = rg.Register(info0)
	require.Nil(t, err)
	time.Sleep(15 * time.Second) // wait register service
	desc := rs.Target(context.TODO(), rpcinfo.NewEndpointInfo(serviceName, "", nil, nil))
	result, err := rs.Resolve(context.TODO(), desc)
	require.Nil(t, err)
	expected := discovery.Result{
		Cacheable: true,
		CacheKey:  serviceName,
		Instances: []discovery.Instance{
			discovery.NewInstance(info0.Addr.Network(), info0.Addr.String(), info0.Weight, info0.Tags),
		},
	}
	require.Equal(t, expected, result)
	Wathcherchange, err := rs.Watcher(context.TODO(), desc)
	t.Logf("the number of instance is %d", len(Wathcherchange.Result.Instances))

	// test register service
	info1 := &registry.Info{
		ServiceName: serviceName,
		Addr:        utils.NewNetAddr("tcp", "127.0.0.1:7777"),
		Weight:      100,
		Tags:        nil,
	}
	err = rg.Register(info1)
	require.Nil(t, err)
	time.Sleep(15 * time.Second) // wait register service
	Wathcherchange, err = rs.Watcher(context.TODO(), desc)
	t.Logf("the number of instance is %d", len(Wathcherchange.Result.Instances))
	result, err = rs.Resolve(context.TODO(), desc)
	require.Nil(t, err)

	// test deregister service

	err = rg.Deregister(info0)
	require.Nil(t, err)
	err = rg.Deregister(info1)
	require.Nil(t, err)
	time.Sleep(15 * time.Second) // wait deregister service
	Wathcherchange, err = rs.Watcher(context.TODO(), desc)
	t.Logf("the number of instance is %d", len(Wathcherchange.Result.Instances))
	desc = rs.Target(context.TODO(), rpcinfo.NewEndpointInfo(serviceName, "", nil, nil))
	result, err = rs.Resolve(context.TODO(), desc)
	require.NotNil(t, err)
}

func TestEmptyEndpoints(t *testing.T) {
	_, err := NewPolarisResolver([]string{})
	require.NotNil(t, err)
}
