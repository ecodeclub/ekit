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

package syncx_test

import (
	"fmt"

	"github.com/ecodeclub/ekit/syncx"
)

func ExampleNewSegmentKeysLock() {
	// 参数就是分多少段，你也可以理解为总共有多少锁
	// 锁越多，并发竞争越低，但是消耗内存；
	// 锁越少，并发竞争越高，但是内存消耗少；
	lock := syncx.NewSegmentKeysLock(100)
	// 对应的还有 TryLock
	// RLock 和 RUnlock
	lock.Lock("key1")
	defer lock.Unlock("key1")
	fmt.Println("OK")
	// Output:
	// OK
}
