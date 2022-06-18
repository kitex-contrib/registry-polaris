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

type Options struct {
	DstNamespace string            `json:"dst_namespace"`
	DstService   string            `json:"dst_service"`
	DstMetadata  map[string]string `json:"dst_metadata"`
	SrcNamespace string            `json:"src_namespace"`
	SrcService   string            `json:"src_service"`
	SrcMetadata  map[string]string `json:"src_metadata"`
}
