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
	"net"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/polarismesh/polaris-go/pkg/model"
)

type polarisKitexInstance struct {
	kitexInstance   discovery.Instance
	polarisInstance model.Instance
	polarisOptions  Options
}

func (i *polarisKitexInstance) Address() net.Addr {
	return i.kitexInstance.Address()
}

func (i *polarisKitexInstance) Weight() int {
	return i.kitexInstance.Weight()
}

func (i *polarisKitexInstance) Tag(key string) (value string, exist bool) {
	value, exist = i.kitexInstance.Tag(key)
	return
}
