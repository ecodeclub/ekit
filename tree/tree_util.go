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
	"reflect"
	"sort"
	"strings"

	"github.com/gotomicro/ekit/bean/option"
)

// Node 树节点的抽象
type Node map[string]any

// Config 配置类
type Config struct {
	IdKey       string // id属性名
	ParentIdKey string // 父节点id属性名
	ChildrenKey string // 子节点属性名
	SortKey     string // 排序属性名
	Deep        int    // 树最大深度(包含根节点) 0:不限制
}

// Builder 存储构造树前的一些元数据,并提供build方法构造树
type Builder[T comparable] struct {
	root Node
	// 标记位,防止build被重复调用
	isBuild bool
	// 配置类
	config *Config
	// Node 的id->Node 的映射
	idTreeMap map[T]Node
	// 按Node 的sort 属性对node的id从小到大进行排序
	sortList []T
}

// NewBuilder 初始化构造器
func NewBuilder[T comparable](rootId T, config *Config) *Builder[T] {
	return &Builder[T]{root: Node{config.IdKey: rootId, config.ChildrenKey: []Node{}}, config: config, isBuild: false}
}

// 对cur增加子节点，同时关联子节点的父节点为当前节点
func addChildren[T comparable](cur Node, builder *Builder[T], children ...Node) {
	if len(children) == 0 {
		return
	}
	config := builder.config

	if m, ok := cur[config.ChildrenKey]; ok {
		cur[config.ChildrenKey] = append(m.([]Node), children...)
	} else {
		cur[config.ChildrenKey] = children
	}
}

func Append[T comparable, E any](builder *Builder[T], list []E, parse func(E) Node) {
	if builder.isBuild {
		panic("树已经被构造")
	}

	if len(builder.idTreeMap) == 0 {
		builder.idTreeMap = make(map[T]Node, len(list))
	}
	if len(builder.sortList) == 0 {
		builder.sortList = make([]T, 0, len(list))
	}

	config := builder.config
	for _, elem := range list {
		node := parse(elem)
		id := node[config.IdKey].(T)
		builder.idTreeMap[id] = node
		builder.sortList = append(builder.sortList, id)
	}

	if config.SortKey != "" {
		sort.Slice(builder.sortList, func(i, j int) bool {
			id1 := builder.sortList[i]
			id2 := builder.sortList[j]
			sort1 := builder.idTreeMap[id1][builder.config.SortKey].(int)
			sort2 := builder.idTreeMap[id2][builder.config.SortKey].(int)
			return sort1 < sort2
		})
	}
}

// Build 构造树,如果原始切片未找到指定的rootTree,会新创建一个Node
func (builder *Builder[T]) Build() Node {
	// 防止重复构造
	if builder.isBuild {
		panic("树已经被构造")
	}
	config := builder.config

	// 按sortKey的顺序遍历Node,因为追加子节点时总是在切片末尾append,所以能保证子节点有序
	for _, id := range builder.sortList {
		node := builder.idTreeMap[id]

		// 遍历到根节点,将根节点上的子节点复制到当前节点,并把当前节点指向根节点
		if id == builder.root[config.IdKey] {
			if children, ok := builder.root[config.ChildrenKey]; ok {
				addChildren(node, builder, children.([]Node)...)
			}
			builder.root = node
			continue
		}
		// 如果父节点是根节点,直接在根节点末尾添加子map
		if parentId := node[config.ParentIdKey]; parentId == builder.root[config.IdKey] {
			addChildren(builder.root, builder, node)
			continue
		}
		// 如果父节点能在map中找到,将当前节点加入
		if parentNode, ok := builder.idTreeMap[node[config.ParentIdKey].(T)]; ok {
			addChildren(parentNode, builder, node)
		}
	}

	if config.Deep > 0 {
		cutTree(builder.root, 1, config.Deep, config)
	}
	builder.isBuild = true
	return builder.root
}

// cutTree 通过递归剪枝
func cutTree(node Node, currentDeep, maxDeeP int, config *Config) {
	if node == nil {
		return
	}
	if currentDeep == maxDeeP {
		delete(node, config.ChildrenKey)
	}
	if children, ok := node[config.ChildrenKey]; ok {
		for _, child := range children.([]Node) {
			cutTree(child, currentDeep+1, maxDeeP, config)
		}
	}
}

func NewConfig(opts ...option.Option[Config]) *Config {
	c := &Config{IdKey: "id", ParentIdKey: "parentId", ChildrenKey: "children", SortKey: "sort"}
	for _, opt := range opts {
		option.Apply(c, opt)
	}
	return c
}
func SetIdKey(key string) option.Option[Config] {
	return func(config *Config) {
		config.IdKey = key
	}
}
func SetParentIdKey(key string) option.Option[Config] {
	return func(config *Config) {
		config.ParentIdKey = key
	}
}

func setChildrenKey(key string) option.Option[Config] {
	return func(config *Config) {
		config.ChildrenKey = key
	}
}
func SetSortKey(key string) option.Option[Config] {
	return func(config *Config) {
		config.SortKey = key
	}
}
func SetDeep(deep int) option.Option[Config] {
	return func(config *Config) {
		config.Deep = deep
	}
}

// GetParser 默认struct->map 转换规则,属性首字母小写,map中的key 必须包含config 中的idKey,parentId Key
// 用户可以自定义一个Config,调用此函数生成一个 Parser 函数,作为参数传递给Append()
func GetParser[T comparable, E any](builder *Builder[T]) func(src E) Node {
	return func(src E) Node {
		return DefaultParser(src, builder)
	}
}

func DefaultParser[T comparable, E any](src E, builder *Builder[T]) Node {
	t, v := reflect.TypeOf(src), reflect.ValueOf(src)
	n := make(Node, t.NumField())
	config := builder.config
	var idFlag, parentFlag bool

	isSort := 1
	if builder.config.SortKey == "" {
		isSort = 0
	}

	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		name := t.Field(i).Name
		key := strings.ToLower(string(name[0])) + name[1:]
		switch key {
		case config.IdKey:
			n[key] = v.Field(i).Interface().(T)
			idFlag = true
		case config.ParentIdKey:
			n[key] = v.Field(i).Interface().(T)
			parentFlag = true
		case config.SortKey:
			n[key] = v.Field(i).Interface().(int)
			isSort = isSort ^ 1
		default:
			n[key] = v.Field(i).Interface()
		}
	}
	if !idFlag {
		panic("未匹配id Key")
	}
	if !parentFlag {
		panic("未匹配parentId Key")
	}
	if isSort == 1 {
		panic("未匹配Sort key")
	}
	return n
}
