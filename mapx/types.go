// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mapx

type mapi[K any, V any] interface {
	Put(key K, val V) error
	Get(key K) (V, bool)
	// Delete 删除
	// 第一个返回值是被删除的 key 对应的值
	// 第二个返回值是代表是否真的删除了
	Delete(k K) (V, bool)
	// Keys 返回所有的键
	// 注意，当你调用多次拿到的结果不一定相等
	// 取决于具体实现
	Keys() []K
	// Values 返回所有的值
	// 注意，当你调用多次拿到的结果不一定相等
	// 取决于具体实现
	Values() []V
	// 返回键值对数量
	Len() int64
}
