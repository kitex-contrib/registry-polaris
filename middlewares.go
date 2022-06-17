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
	"time"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/remote/codec/perrors"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/polarismesh/polaris-go/api"
)

func NewUpdateServiceCallResultMW(configFile ...string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request, response interface{}) error {
			retCode := int32(-1) // todo 成功的code？
			retStatus := api.RetSuccess
			begin := time.Now()
			kitexCallErr := next(ctx, request, response)
			cost := time.Since(begin)
			if kitexCallErr != nil {
				if e, ok := kitexCallErr.(perrors.ProtocolError); ok {
					retCode = int32(e.TypeId())
				} else {
					retCode = perrors.UnknownProtocolError
				}
				retStatus = api.RetFail
			}

			ri := rpcinfo.GetRPCInfo(ctx)
			sdkCtx, err := GetPolarisConfig(configFile...)
			if err != nil {
				return err
			}
			consumer := api.NewConsumerAPIByContext(sdkCtx)
			ns, _ := ri.To().Tag("namespace")
			instanceId, ok := ri.To().Tag("instanceId")
			if !ok {
				// 没有找到实例
				return kitexCallErr
			}

			// todo del
			fmt.Printf("call addr=%s \n", ri.To().Address().String())

			req := api.InstanceRequest{
				ServiceKey: model.ServiceKey{
					Namespace: ns,
					Service:   ri.To().ServiceName(),
				},
				InstanceID: instanceId,
			}

			svcCallResult, err := api.NewServiceCallResult(sdkCtx, req)
			if err != nil {
				return err
			}

			svcCallResult.SetRetCode(retCode)
			svcCallResult.SetRetStatus(retStatus)
			svcCallResult.SetDelay(cost)
			// 执行调用结果上报
			_ = consumer.UpdateServiceCallResult(svcCallResult)
			return kitexCallErr
		}
	}
}
