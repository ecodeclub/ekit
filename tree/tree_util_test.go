// Copyright 2021 gotomicro
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

package tree

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Area struct {
	Id       int
	ParentId int
	Name     string
	Sort     int
}

func MockAreaList() []Area {
	return []Area{
		{1, -1, "中国", 1},
		{2, 1, "北京", 4},
		{3, 1, "上海", 3},
		{4, 2, "朝阳区", 5},
		{5, 2, "海淀区", 6},
		{6, 3, "浦东", 6},
		{7, 3, "长宁", 1},
	}
}

func ExampleBuilder_Build() {
	builder := NewBuilder(1, NewConfig())
	Append(builder, MockAreaList()[:5], GetParser[int, Area](builder))
	mp := builder.Build()
	str, _ := json.Marshal(mp)
	fmt.Println(string(str))
	// Output:
	// {"children":[{"id":3,"name":"上海","parentId":1,"sort":3},{"children":[{"id":4,"name":"朝阳区","parentId":2,"sort":5},{"id":5,"name":"海淀区","parentId":2,"sort":6}],"id":2,"name":"北京","parentId":1,"sort":4}],"id":1,"name":"中国","parentId":-1,"sort":1}
}

func TestBuilder_Build(t *testing.T) {

	tests := []struct {
		name   string
		config *Config
		list   []Area
		parse  func(src Area) Node
		want   map[string]any
		rootId int
	}{
		{
			name:   "单节点树构造",
			config: NewConfig(),
			list:   MockAreaList()[:1],
			parse:  nil,
			want:   Node{"id": 1, "parentId": -1, "name": "中国", "sort": 1},
			rootId: 1,
		},
		{
			name:   "设置树深度",
			config: NewConfig(SetDeep(2)),
			list:   MockAreaList(),
			parse:  nil,
			want:   Node{"id": 1, "parentId": -1, "name": "中国", "sort": 1, "children": []Node{{"id": 3, "parentId": 1, "name": "上海", "sort": 3}, {"id": 2, "parentId": 1, "name": "北京", "sort": 4}}},
			rootId: 1,
		},
		{
			name:   "自定义mapKey",
			config: NewConfig(SetDeep(2), SetIdKey("code"), SetParentIdKey("parentCode"), setChildrenKey("child's")),
			list:   MockAreaList(),
			parse: func(src Area) Node {
				return Node{"code": src.Id, "parentCode": src.ParentId, "sort": src.Sort, "name": src.Name}
			},
			want:   Node{"code": 1, "parentCode": -1, "name": "中国", "sort": 1, "child's": []Node{{"code": 3, "parentCode": 1, "name": "上海", "sort": 3}, {"code": 2, "parentCode": 1, "name": "北京", "sort": 4}}},
			rootId: 1,
		}, {
			name:   "取消排序策略",
			config: NewConfig(SetSortKey("")),
			list:   MockAreaList()[:3],
			parse:  nil,
			want:   Node{"id": 1, "parentId": -1, "name": "中国", "sort": 1, "children": []Node{{"id": 2, "parentId": 1, "name": "北京", "sort": 4}, {"id": 3, "parentId": 1, "name": "上海", "sort": 3}}},
			rootId: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(tt.rootId, tt.config)
			if tt.parse == nil {
				tt.parse = GetParser[int, Area](builder)
			}
			Append(builder, tt.list, tt.parse)
			got := builder.Build()
			json1, _ := json.Marshal(tt.want)
			json2, _ := json.Marshal(got)
			assert.Equal(t, string(json1), string(json2))
		})
	}
}
