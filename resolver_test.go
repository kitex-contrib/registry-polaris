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
	serviceName = "registry-test"
)

func TestPolarisResolver(t *testing.T) {
	rg, err := NewPolarisRegistry()
	require.Nil(t, err)
	rs, err := NewPolarisResolver()
	require.Nil(t, err)

	// test register service
	InstanceOne := &registry.Info{
		ServiceName: serviceName,
		Addr:        utils.NewNetAddr("tcp", "127.0.0.1:6666"),
		Weight:      100,
		Tags:        nil, // when Tags is nil the namespace is default
	}
	err = rg.Register(InstanceOne)
	require.Nil(t, err)
	time.Sleep(15 * time.Second)                                                          // wait register service
	desc := rs.Target(context.TODO(), rpcinfo.NewEndpointInfo(serviceName, "", nil, nil)) // the namespace is default
	result, err := rs.Resolve(context.TODO(), desc)
	require.Nil(t, err)
	expected := discovery.Result{
		Cacheable: true,
		CacheKey:  polarisDefaultNamespace + ":" + serviceName,
		Instances: []discovery.Instance{
			discovery.NewInstance(InstanceOne.Addr.Network(), InstanceOne.Addr.String(), InstanceOne.Weight, map[string]string{
				"namespace": "default",
			}),
		},
	}
	require.Equal(t, expected, result)
	watcherChange, err := rs.Watcher(context.TODO(), desc)
	require.Nil(t, err)
	t.Logf("the number of instance is %d", len(watcherChange.Result.Instances))

	// test register service
	InstanceTwo := &registry.Info{
		ServiceName: serviceName,
		Addr:        utils.NewNetAddr("tcp", "127.0.0.1:7777"),
		Weight:      100,
		Tags:        nil, // namespace is default
	}
	err = rg.Register(InstanceTwo)
	require.Nil(t, err)
	time.Sleep(15 * time.Second) // wait register service
	watcherChange, err = rs.Watcher(context.TODO(), desc)
	require.Nil(t, err)
	t.Logf("the number of instance is %d", len(watcherChange.Result.Instances))
	result, err = rs.Resolve(context.TODO(), desc)
	require.Nil(t, err)

	// test deregister service

	err = rg.Deregister(InstanceOne) // deregister InstanceOne
	require.Nil(t, err)
	err = rg.Deregister(InstanceTwo) // deregister InstanceTwo
	require.Nil(t, err)
	time.Sleep(15 * time.Second) // wait deregister service
	watcherChange, err = rs.Watcher(context.TODO(), desc)
	require.Nil(t, err)
	t.Logf("the number of instance is %d", len(watcherChange.Result.Instances))
	desc = rs.Target(context.TODO(), rpcinfo.NewEndpointInfo(serviceName, "", nil, nil)) // namespace is  default
	result, err = rs.Resolve(context.TODO(), desc)
	require.NotNil(t, err)
}

func TestEmptyEndpoints(t *testing.T) {
	_, err := NewPolarisResolver()
	require.Nil(t, err)
}
