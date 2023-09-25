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

package copier

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ecodeclub/ekit/bean/copier/converter"

	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

func TestReflectCopier_Copy(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		copyFunc func() (any, error)
		wantDst  any
		wantErr  error
	}{
		{
			name: "simple struct",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				})
			},
			wantDst: &SimpleDst{
				Name:    "大明",
				Age:     ekit.ToPtr[int](18),
				Friends: []string{"Tom", "Jerry"},
			},
		},
		{
			name: "基础类型的 struct",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[BasicSrc, BasicDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&BasicSrc{
					Name:    "大明",
					Age:     10,
					CNumber: complex(1, 2),
				})
			},
			wantDst: &BasicDst{
				Name:    "大明",
				Age:     10,
				CNumber: complex(1, 2),
			},
		},
		{
			name: "src 是基础类型",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[int, int]()
				if err != nil {
					return nil, err
				}
				i := 10
				return copier.Copy(&i)
			},
			wantErr: newErrTypeError(reflect.TypeOf(10)),
		},
		{
			name: "dst 是基础类型",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, string]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				})
			},
			wantErr: newErrTypeError(reflect.TypeOf("")),
		},
		{
			name: "接口类型",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[InterfaceSrc, InterfaceDst]()
				if err != nil {
					return nil, err
				}
				i := InterfaceSrc(10)
				return copier.Copy(&i)
			},
			wantErr: newErrTypeError(reflect.TypeOf(new(InterfaceSrc)).Elem()),
		},
		{
			name: "simple struct 空切片, 空指针",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name: "大明",
				})
			},
			wantDst: &SimpleDst{
				Name: "大明",
			},
		},
		{
			name: "组合 struct ",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[EmbedSrc, EmbedDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&EmbedSrc{
					SimpleSrc: SimpleSrc{
						Name:    "xiaoli",
						Age:     ekit.ToPtr[int](19),
						Friends: []string{},
					},
					BasicSrc: &BasicSrc{
						Name:    "xiaowang",
						Age:     20,
						CNumber: complex(2, 2),
					},
				})
			},
			wantDst: &EmbedDst{
				SimpleSrc: SimpleSrc{
					Name:    "xiaoli",
					Age:     ekit.ToPtr[int](19),
					Friends: []string{},
				},
				BasicSrc: &BasicSrc{
					Name:    "xiaowang",
					Age:     20,
					CNumber: complex(2, 2),
				},
			},
		},
		{
			name: "复杂 Struct",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ComplexSrc, ComplexDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ComplexSrc{
					Simple: SimpleSrc{
						Name:    "xiaohong",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrc: SimpleSrc{
							Name:    "xiaopeng",
							Age:     ekit.ToPtr[int](88),
							Friends: []string{"la", "ha", "le"},
						},
						BasicSrc: &BasicSrc{
							Name:    "wang",
							Age:     22,
							CNumber: complex(2, 1),
						},
					},
					BasicSrc: BasicSrc{
						Name:    "wang11",
						Age:     22,
						CNumber: complex(2, 1),
					},
				})
			},
			wantDst: &ComplexDst{
				Simple: SimpleDst{
					Name:    "xiaohong",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"ha", "ha", "le"},
				},
				Embed: &EmbedDst{
					SimpleSrc: SimpleSrc{
						Name:    "xiaopeng",
						Age:     ekit.ToPtr[int](88),
						Friends: []string{"la", "ha", "le"},
					},
					BasicSrc: &BasicSrc{
						Name:    "wang",
						Age:     22,
						CNumber: complex(2, 1),
					},
				},
				BasicSrc: BasicSrc{
					Name:    "wang11",
					Age:     22,
					CNumber: complex(2, 1),
				},
			},
		},
		{
			name: "特殊类型",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SpecialSrc, SpecialDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SpecialSrc{
					Arr: [3]float32{1, 2, 3},
					M: map[string]int{
						"ha": 1,
						"o":  2,
					},
				})
			},
			wantDst: &SpecialDst{
				Arr: [3]float32{1, 2, 3},
				M: map[string]int{
					"ha": 1,
					"o":  2,
				},
			},
		},
		{
			name: "复杂 Struct 不匹配",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[NotMatchSrc, NotMatchDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&NotMatchSrc{
					Simple: SimpleSrc{
						Name:    "xiaohong",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrc: SimpleSrc{
							Name:    "xiaopeng",
							Age:     ekit.ToPtr[int](88),
							Friends: []string{"la", "ha", "le"},
						},
						BasicSrc: &BasicSrc{
							Name:    "wang",
							Age:     22,
							CNumber: complex(2, 1),
						},
					},
					BasicSrc: BasicSrc{
						Name:    "wang11",
						Age:     22,
						CNumber: complex(2, 1),
					},
					S: struct{ A string }{A: "a"},
				})
			},
			wantErr: newErrTypeNotMatchError(reflect.TypeOf(""), reflect.TypeOf(0), "A"),
		},
		{
			name: "多重指针",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[MultiPtrSrc, MultiPtrDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&MultiPtrSrc{
					Name:    "a",
					Age:     ekit.ToPtr[*int](ekit.ToPtr[int](10)),
					Friends: nil,
				})
			},
			wantErr: newErrMultiPointer("Age"),
		},
		{
			name: "src 有额外字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[DiffSrc, DiffDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&DiffSrc{
					A: "xiaowang",
					B: 100,
					c: SimpleSrc{
						Name: "66",
						Age:  ekit.ToPtr[int](100),
					},
					F: BasicSrc{
						Name:    "good name",
						Age:     200,
						CNumber: complex(2, 2),
					},
				})
			},
			wantDst: &DiffDst{
				A: "xiaowang",
				B: 100,
				d: SimpleSrc{},
				G: BasicSrc{},
			},
		},
		{
			name: "dst 有额外字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[DiffSrc, DiffDst]()
				if err != nil {
					return nil, err
				}
				dst := &DiffDst{
					A: "66",
					B: 1,
					d: SimpleSrc{
						Name: "wodemingzi",
						Age:  ekit.ToPtr(int(10)),
					},
					G: BasicSrc{
						Name:    "nidemingzi",
						Age:     23,
						CNumber: complex(1, 2),
					},
				}
				err = copier.CopyTo(&DiffSrc{
					A: "xiaowang",
					B: 100,
					c: SimpleSrc{
						Name: "66",
						Age:  ekit.ToPtr[int](100),
					},
					F: BasicSrc{
						Name:    "good name",
						Age:     200,
						CNumber: complex(2, 2),
					},
				}, dst)
				return dst, err
			},
			wantDst: &DiffDst{
				A: "xiaowang",
				B: 100,
				d: SimpleSrc{
					Name: "wodemingzi",
					Age:  ekit.ToPtr(int(10)),
				},
				G: BasicSrc{
					Name:    "nidemingzi",
					Age:     23,
					CNumber: complex(1, 2),
				},
			},
		},
		{
			name: "跨层级别匹配",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleEmbedDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{})
			},
			wantDst: &SimpleEmbedDst{},
		},
		{
			name: "成员为结构体数组",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ArraySrc, ArrayDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ArraySrc{
					A: []SimpleSrc{
						{
							Name:    "大明",
							Age:     ekit.ToPtr[int](18),
							Friends: []string{"Tom", "Jerry"},
						},
						{
							Name:    "小明",
							Age:     ekit.ToPtr[int](8),
							Friends: []string{"Tom"},
						},
					},
				})
			},
			wantDst: &ArrayDst{
				A: []SimpleSrc{
					{
						Name:    "大明",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"Tom", "Jerry"},
					},
					{
						Name:    "小明",
						Age:     ekit.ToPtr[int](8),
						Friends: []string{"Tom"},
					},
				},
			},
		},
		{
			name: "成员为结构体数组，结构体不同",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ArraySrc, ArrayDst1]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ArraySrc{
					A: []SimpleSrc{
						{
							Name:    "大明",
							Age:     ekit.ToPtr[int](18),
							Friends: []string{"Tom", "Jerry"},
						},
						{
							Name:    "小明",
							Age:     ekit.ToPtr[int](8),
							Friends: []string{"Tom"},
						},
					},
				})
			},
			wantErr: newErrTypeNotMatchError(reflect.TypeOf(new([]SimpleSrc)).Elem(), reflect.TypeOf(new([]SimpleDst)).Elem(), "A"),
		},
		{
			name: "成员为map结构体",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[MapSrc, MapDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&MapSrc{
					A: map[string]SimpleSrc{
						"a": {
							Name:    "大明",
							Age:     ekit.ToPtr[int](18),
							Friends: []string{"Tom", "Jerry"},
						},
					},
				})
			},
			wantDst: &MapDst{
				A: map[string]SimpleSrc{
					"a": {
						Name:    "大明",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"Tom", "Jerry"},
					},
				},
			},
		},
		{
			name: "成员为不同的map结构体",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[MapSrc, MapDst1]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&MapSrc{
					A: map[string]SimpleSrc{
						"a": {
							Name:    "大明",
							Age:     ekit.ToPtr[int](18),
							Friends: []string{"Tom", "Jerry"},
						},
					},
				})
			},
			wantErr: newErrTypeNotMatchError(reflect.TypeOf(new(map[string]SimpleSrc)).Elem(), reflect.TypeOf(new(map[string]SimpleDst)).Elem(), "A"),
		},
		{
			name: "成员有别名类型",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SpecialSrc1, SpecialDst1]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SpecialSrc1{
					A: 1,
				})
			},
			wantErr: newErrTypeNotMatchError(reflect.TypeOf(new(int)).Elem(), reflect.TypeOf(new(aliasInt)).Elem(), "A"),
		},
		{
			name: "成员有别名类型1",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SpecialSrc1, SpecialDst2]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SpecialSrc1{
					A: 1,
				})
			},
			wantDst: &SpecialDst2{A: 1},
		},
		{
			name: "simple_struct_忽略字段的时候传空",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				}, IgnoreFields())
			},
			wantDst: &SimpleDst{
				Name:    "大明",
				Age:     ekit.ToPtr[int](18),
				Friends: []string{"Tom", "Jerry"},
			},
		},
		{
			name: "simple_struct_忽略一个字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				}, IgnoreFields("Age"))
			},
			wantDst: &SimpleDst{
				Name:    "大明",
				Age:     nil,
				Friends: []string{"Tom", "Jerry"},
			},
		},
		{
			name: "simple_struct_忽略多个字段_传入多个Option_每个Option传入一个字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				}, IgnoreFields("Age"), IgnoreFields("Friends"))
			},
			wantDst: &SimpleDst{
				Name:    "大明",
				Age:     nil,
				Friends: nil,
			},
		},
		{
			name: "simple_struct_忽略多个字段_传入一个Option_Option传入多个字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				}, IgnoreFields("Age", "Friends"))
			},
			wantDst: &SimpleDst{
				Name:    "大明",
				Age:     nil,
				Friends: nil,
			},
		},
		{
			name: "simple_struct_忽略全部字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				}, IgnoreFields("Name"), IgnoreFields("Age"), IgnoreFields("Friends"))
			},
			wantDst: &SimpleDst{},
		},
		{
			name: "simple_struct_空切片_空指针_忽略字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrc{
					Name: "大明",
				}, IgnoreFields("Name"))
			},
			wantDst: &SimpleDst{
				Name: "",
			},
		},
		{
			name: "组合_struct_忽略组合中的一个字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[EmbedSrc, EmbedDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&EmbedSrc{
					SimpleSrc: SimpleSrc{
						Name:    "xiaoli",
						Age:     ekit.ToPtr[int](19),
						Friends: []string{},
					},
					BasicSrc: &BasicSrc{
						Name:    "xiaowang",
						Age:     20,
						CNumber: complex(2, 2),
					},
				}, IgnoreFields("CNumber"))
			},
			wantDst: &EmbedDst{
				SimpleSrc: SimpleSrc{
					Name:    "xiaoli",
					Age:     ekit.ToPtr[int](19),
					Friends: []string{},
				},
				BasicSrc: &BasicSrc{
					Name:    "xiaowang",
					Age:     20,
					CNumber: complex(0, 0),
				},
			},
		},
		{
			name: "组合_struct_忽略组合中全部同名字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[EmbedSrc, EmbedDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&EmbedSrc{
					SimpleSrc: SimpleSrc{
						Name:    "xiaoli",
						Age:     ekit.ToPtr[int](19),
						Friends: []string{},
					},
					BasicSrc: &BasicSrc{
						Name:    "xiaowang",
						Age:     20,
						CNumber: complex(2, 2),
					},
				}, IgnoreFields("Age"))
			},
			wantDst: &EmbedDst{
				SimpleSrc: SimpleSrc{
					Name:    "xiaoli",
					Age:     nil,
					Friends: []string{},
				},
				BasicSrc: &BasicSrc{
					Name:    "xiaowang",
					Age:     0,
					CNumber: complex(2, 2),
				},
			},
		},
		{
			name: "组合_struct_忽略组合中同名结构体",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[EmbedSrc, EmbedDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&EmbedSrc{
					SimpleSrc: SimpleSrc{
						Name:    "xiaoli",
						Age:     ekit.ToPtr[int](19),
						Friends: []string{},
					},
					BasicSrc: &BasicSrc{
						Name:    "xiaowang",
						Age:     20,
						CNumber: complex(2, 2),
					},
				}, IgnoreFields("SimpleSrc"))
			},
			wantDst: &EmbedDst{
				SimpleSrc: SimpleSrc{},
				BasicSrc: &BasicSrc{
					Name:    "xiaowang",
					Age:     20,
					CNumber: complex(2, 2),
				},
			},
		},
		{
			name: "复杂_Struct_忽略多层嵌套中全部同名字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ComplexSrc, ComplexDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ComplexSrc{
					Simple: SimpleSrc{
						Name:    "xiaohong",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrc: SimpleSrc{
							Name:    "xiaopeng",
							Age:     ekit.ToPtr[int](88),
							Friends: []string{"la", "ha", "le"},
						},
						BasicSrc: &BasicSrc{
							Name:    "wang",
							Age:     22,
							CNumber: complex(2, 1),
						},
					},
					BasicSrc: BasicSrc{
						Name:    "wang11",
						Age:     22,
						CNumber: complex(2, 1),
					},
				}, IgnoreFields("Age"))
			},
			wantDst: &ComplexDst{
				Simple: SimpleDst{
					Name:    "xiaohong",
					Age:     nil,
					Friends: []string{"ha", "ha", "le"},
				},
				Embed: &EmbedDst{
					SimpleSrc: SimpleSrc{
						Name:    "xiaopeng",
						Age:     nil,
						Friends: []string{"la", "ha", "le"},
					},
					BasicSrc: &BasicSrc{
						Name:    "wang",
						Age:     0,
						CNumber: complex(2, 1),
					},
				},
				BasicSrc: BasicSrc{
					Name:    "wang11",
					Age:     0,
					CNumber: complex(2, 1),
				},
			},
		},
		{
			name: "复杂_Struct_忽略多层嵌套中的同名结构体",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ComplexSrc, ComplexDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ComplexSrc{
					Simple: SimpleSrc{
						Name:    "xiaohong",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrc: SimpleSrc{
							Name:    "xiaopeng",
							Age:     ekit.ToPtr[int](88),
							Friends: []string{"la", "ha", "le"},
						},
						BasicSrc: &BasicSrc{
							Name:    "wang",
							Age:     22,
							CNumber: complex(2, 1),
						},
					},
					BasicSrc: BasicSrc{
						Name:    "wang11",
						Age:     22,
						CNumber: complex(2, 1),
					},
				}, IgnoreFields("SimpleSrc"))
			},
			wantDst: &ComplexDst{
				Simple: SimpleDst{
					Name:    "xiaohong",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"ha", "ha", "le"},
				},
				Embed: &EmbedDst{
					SimpleSrc: SimpleSrc{},
					BasicSrc: &BasicSrc{
						Name:    "wang",
						Age:     22,
						CNumber: complex(2, 1),
					},
				},
				BasicSrc: BasicSrc{
					Name:    "wang11",
					Age:     22,
					CNumber: complex(2, 1),
				},
			},
		},
		{
			name: "复杂_Struct_忽略多层嵌套中的整个结构体",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ComplexSrc, ComplexDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ComplexSrc{
					Simple: SimpleSrc{
						Name:    "xiaohong",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrc: SimpleSrc{
							Name:    "xiaopeng",
							Age:     ekit.ToPtr[int](88),
							Friends: []string{"la", "ha", "le"},
						},
						BasicSrc: &BasicSrc{
							Name:    "wang",
							Age:     22,
							CNumber: complex(2, 1),
						},
					},
					BasicSrc: BasicSrc{
						Name:    "wang11",
						Age:     22,
						CNumber: complex(2, 1),
					},
				}, IgnoreFields("Embed"))
			},
			wantDst: &ComplexDst{
				Simple: SimpleDst{
					Name:    "xiaohong",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"ha", "ha", "le"},
				},
				Embed: nil,
				BasicSrc: BasicSrc{
					Name:    "wang11",
					Age:     22,
					CNumber: complex(2, 1),
				},
			},
		},
		{
			name: "特殊类型_忽略结构体中的切片",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SpecialSrc, SpecialDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SpecialSrc{
					Arr: [3]float32{1, 2, 3},
					M: map[string]int{
						"ha": 1,
						"o":  2,
					},
				}, IgnoreFields("Arr"))
			},
			wantDst: &SpecialDst{
				Arr: [3]float32{},
				M: map[string]int{
					"ha": 1,
					"o":  2,
				},
			},
		},
		{
			name: "特殊类型_忽略结构体中的map",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SpecialSrc, SpecialDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SpecialSrc{
					Arr: [3]float32{1, 2, 3},
					M: map[string]int{
						"ha": 1,
						"o":  2,
					},
				}, IgnoreFields("M"))
			},
			wantDst: &SpecialDst{
				Arr: [3]float32{1, 2, 3},
				M:   nil,
			},
		},
		{
			name: "dst_有额外字段_忽略一个字段_其他字段会被赋值",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[DiffSrc, DiffDst]()
				if err != nil {
					return nil, err
				}
				dst := &DiffDst{
					A: "66",
					B: 1,
					d: SimpleSrc{
						Name: "wodemingzi",
						Age:  ekit.ToPtr(int(10)),
					},
					G: BasicSrc{
						Name:    "nidemingzi",
						Age:     23,
						CNumber: complex(1, 2),
					},
				}
				err = copier.CopyTo(&DiffSrc{
					A: "xiaowang",
					B: 100,
					c: SimpleSrc{
						Name: "66",
						Age:  ekit.ToPtr[int](100),
					},
					F: BasicSrc{
						Name:    "good name",
						Age:     200,
						CNumber: complex(2, 2),
					},
				}, dst, IgnoreFields("A"))
				return dst, err
			},
			wantDst: &DiffDst{
				A: "66",
				B: 100,
				d: SimpleSrc{
					Name: "wodemingzi",
					Age:  ekit.ToPtr(int(10)),
				},
				G: BasicSrc{
					Name:    "nidemingzi",
					Age:     23,
					CNumber: complex(1, 2),
				},
			},
		},
		{
			name: "dst_有额外字段_不会忽略dst的字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[DiffSrc, DiffDst]()
				if err != nil {
					return nil, err
				}
				dst := &DiffDst{
					A: "66",
					B: 1,
					d: SimpleSrc{
						Name: "wodemingzi",
						Age:  ekit.ToPtr(int(10)),
					},
					G: BasicSrc{
						Name:    "nidemingzi",
						Age:     23,
						CNumber: complex(1, 2),
					},
				}
				err = copier.CopyTo(&DiffSrc{
					A: "xiaowang",
					B: 100,
					c: SimpleSrc{
						Name: "66",
						Age:  ekit.ToPtr[int](100),
					},
					F: BasicSrc{
						Name:    "good name",
						Age:     200,
						CNumber: complex(2, 2),
					},
				}, dst, IgnoreFields("G"))
				return dst, err
			},
			wantDst: &DiffDst{
				A: "xiaowang",
				B: 100,
				d: SimpleSrc{
					Name: "wodemingzi",
					Age:  ekit.ToPtr(int(10)),
				},
				G: BasicSrc{
					Name:    "nidemingzi",
					Age:     23,
					CNumber: complex(1, 2),
				},
			},
		},
		{
			name: "成员为结构体数组_不会忽略结构体中的字段",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ArraySrc, ArrayDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ArraySrc{
					A: []SimpleSrc{
						{
							Name:    "大明",
							Age:     ekit.ToPtr[int](18),
							Friends: []string{"Tom", "Jerry"},
						},
						{
							Name:    "小明",
							Age:     ekit.ToPtr[int](8),
							Friends: []string{"Tom"},
						},
					},
				}, IgnoreFields("Age"))
			},
			wantDst: &ArrayDst{
				A: []SimpleSrc{
					{
						Name:    "大明",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"Tom", "Jerry"},
					},
					{
						Name:    "小明",
						Age:     ekit.ToPtr[int](8),
						Friends: []string{"Tom"},
					},
				},
			},
		},
		{
			name: "指定convert time2string,src为nil",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ConvSimpleSrc{}, ConvertField[time.Time, string]("BirthDay", converter.Time2String{Pattern: "2006-01-02 15:04:05"}))
			},
			wantDst: &ConvSimpleDst{
				BirthDay: "0001-01-01 00:00:00",
			},
		},
		{
			name: "指定convert time2string",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ConvSimpleSrc{
					Name:     "大明",
					BirthDay: time.Date(2023, time.July, 26, 9, 15, 22, 213, time.UTC),
					Friends:  []string{"Tom", "Jerry"},
				}, ConvertField[time.Time, string]("BirthDay", converter.Time2String{Pattern: "2006-01-02 15:04:05"}))
			},
			wantDst: &ConvSimpleDst{
				Name:     "大明",
				BirthDay: "2023-07-26 09:15:22",
				Friends:  []string{"Tom", "Jerry"},
			},
		},
		{
			name: "指定convert func, src为nil",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(
					&ConvSimpleSrc{},
					ConvertField[string, string](
						"Name",
						converter.ConverterFunc[string, string](func(src string) (string, error) {
							newS := fmt.Sprintf("%s plus", src)
							return newS, nil
						}),
					),
					ConvertField[time.Time, string](
						"BirthDay",
						converter.ConverterFunc[time.Time, string](func(src time.Time) (string, error) {
							return src.Format("2006-01-02 15:04:05"), nil
						}),
					),
					ConvertField[*int, *int](
						"Age",
						converter.ConverterFunc[*int, *int](func(src *int) (*int, error) {
							newS := *src + 1
							return &newS, nil
						}),
					),
					ConvertField[[]string, []string](
						"Friends",
						converter.ConverterFunc[[]string, []string](func(src []string) ([]string, error) {
							return []string{"Tom", "Jerry"}, nil
						}),
					),
				)
			},
			wantDst: &ConvSimpleDst{
				Name:     " plus",
				Age:      nil,
				BirthDay: "0001-01-01 00:00:00",
				Friends:  []string{"Tom", "Jerry"},
			},
		},
		{
			name: "指定convert func, dst值为nil",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(
					&ConvSimpleSrc{
						Name:     "大明",
						Age:      ekit.ToPtr[int](11),
						BirthDay: time.Now(),
						Friends:  []string{"Tom", "Jerry"},
					},
					ConvertField[string, string](
						"Name",
						converter.ConverterFunc[string, string](func(src string) (string, error) {
							return "", nil
						}),
					),
					ConvertField[time.Time, string](
						"BirthDay",
						converter.ConverterFunc[time.Time, string](func(src time.Time) (string, error) {
							return "", nil
						}),
					),
					ConvertField[*int, *int](
						"Age",
						converter.ConverterFunc[*int, *int](func(src *int) (*int, error) {
							return nil, nil
						}),
					),
					ConvertField[[]string, []string](
						"Friends",
						converter.ConverterFunc[[]string, []string](func(src []string) ([]string, error) {
							return nil, nil
						}),
					),
				)
			},
			wantDst: &ConvSimpleDst{
				Name:     "",
				BirthDay: "",
				Age:      nil,
				Friends:  nil,
			},
		},
		{
			name: "指定convert func",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(
					&ConvSimpleSrc{
						Name:     "大明",
						Age:      ekit.ToPtr[int](15),
						BirthDay: time.Date(2023, time.July, 26, 9, 15, 22, 213, time.UTC),
						Friends:  []string{"Tom", "Jerry"},
					},
					ConvertField[string, string](
						"Name",
						converter.ConverterFunc[string, string](func(src string) (string, error) {
							newS := fmt.Sprintf("%s plus", src)
							return newS, nil
						}),
					),
					ConvertField[time.Time, string](
						"BirthDay",
						converter.ConverterFunc[time.Time, string](func(src time.Time) (string, error) {
							return src.Format("2006-01-02 15:04:05"), nil
						}),
					),
					ConvertField[*int, *int](
						"Age",
						converter.ConverterFunc[*int, *int](func(src *int) (*int, error) {
							newS := *src + 1
							return &newS, nil
						}),
					),
				)
			},
			wantDst: &ConvSimpleDst{
				Name:     "大明 plus",
				Age:      ekit.ToPtr[int](16),
				BirthDay: "2023-07-26 09:15:22",
				Friends:  []string{"Tom", "Jerry"},
			},
		},
		{
			name: "指定返回特殊类型的convert func",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSpecialSrc, ConvSpecialDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ConvSpecialSrc{
					Arr:  [3]float32{1, 2, 3},
					M:    map[string]int{"a": 4, "b": 5, "c": 6},
					Diff: map[string]int{"a1": 41, "b1": 51, "c1": 61},
				}, ConvertField[map[string]int, map[string]int](
					"M",
					converter.ConverterFunc[map[string]int, map[string]int](func(src map[string]int) (map[string]int, error) {
						newM := map[string]int{"a1": 41, "b1": 51, "c1": 61}
						return newM, nil
					})),
					ConvertField[map[string]int, []int](
						"Diff",
						converter.ConverterFunc[map[string]int, []int](func(src map[string]int) ([]int, error) {
							newM := []int{1, 1, 1}
							return newM, nil
						})),
				)
			},
			wantDst: &ConvSpecialDst{
				Arr:  [3]float32{1, 2, 3},
				M:    map[string]int{"a1": 41, "b1": 51, "c1": 61},
				Diff: []int{1, 1, 1},
			},
		},
		{
			name: "创建时指定默认converter",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst](
					ConvertField[time.Time, string](
						"BirthDay",
						converter.Time2String{Pattern: "2006-01-02 15:04:05"},
					),
				)
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ConvSimpleSrc{
					Name:     "大明",
					BirthDay: time.Date(2023, time.July, 26, 9, 15, 22, 213, time.UTC),
					Friends:  []string{"Tom", "Jerry"},
				}, ConvertField[string, string]("Name", converter.ConverterFunc[string, string](func(src string) (string, error) {
					newS := fmt.Sprintf("%s plus", src)
					return newS, nil
				})))
			},
			wantDst: &ConvSimpleDst{
				Name:     "大明 plus",
				BirthDay: "2023-07-26 09:15:22",
				Friends:  []string{"Tom", "Jerry"},
			},
		},
		{
			name: "创建时指定默认converter,convert同一个字段会覆盖",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst](
					ConvertField[time.Time, string](
						"BirthDay",
						converter.Time2String{Pattern: "2006-01-02 15:04:05"},
					),
				)
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ConvSimpleSrc{
					BirthDay: time.Date(2023, time.July, 26, 9, 15, 22, 213, time.UTC),
				}, ConvertField[time.Time, string]("BirthDay", converter.ConverterFunc[time.Time, string](func(src time.Time) (string, error) {
					return "1234567", nil
				})))
			},
			wantDst: &ConvSimpleDst{
				BirthDay: "1234567",
			},
		},
		{
			name: "创建时指定默认converter,convert同一个字段会覆盖,覆盖后不影响默认配置",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst](
					ConvertField[time.Time, string](
						"BirthDay",
						converter.Time2String{Pattern: "2006-01-02 15:04:05"},
					),
				)
				if err != nil {
					return nil, err
				}
				// 第一次执行Copy,函数中指定converter
				_, err = copier.Copy(
					&ConvSimpleSrc{BirthDay: time.Date(2023, time.July, 26, 9, 15, 22, 213, time.UTC)},
					ConvertField[time.Time, string](
						"BirthDay",
						converter.ConverterFunc[time.Time, string](func(src time.Time) (string, error) {
							return "1234567", nil
						})))
				if err != nil {
					return nil, err
				}
				// 第二次执行Copy,函数中不指定converter,走默认
				return copier.Copy(&ConvSimpleSrc{
					BirthDay: time.Date(2023, time.July, 26, 9, 15, 22, 213, time.UTC),
				})
			},
			wantDst: &ConvSimpleDst{
				BirthDay: "2023-07-26 09:15:22",
			},
		},
		{
			name: "创建时指定默认忽略字段,Copy()时指定的忽略字段不影响默认",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[SimpleSrc, SimpleDst](IgnoreFields("Age"))
				if err != nil {
					return nil, err
				}
				// 第一次执行Copy,函数中指定ignore字段
				_, err = copier.Copy(&SimpleSrc{
					Name: "大明",
					Age:  ekit.ToPtr[int](11),
				}, IgnoreFields("Name"))
				if err != nil {
					return nil, err
				}
				// 第二次执行Copy,函数中不指定ignore字段,走默认
				return copier.Copy(&SimpleSrc{
					Name: "大明",
				})
			},
			wantDst: &SimpleDst{
				Name: "大明",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.copyFunc()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantDst, res)
		})
	}
}

func Test_Concurrency_Copy(t *testing.T) {
	copier, err := NewReflectCopier[ConvSimpleSrc, ConvSimpleDst](
		ConvertField[time.Time, string](
			"BirthDay",
			converter.Time2String{Pattern: "2006-01-02 15:04:05"},
		),
	)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			val := strconv.Itoa(i)
			c, err := copier.Copy(
				&ConvSimpleSrc{BirthDay: time.Date(2023, time.July, 26, 9, 15, 22, 213, time.UTC)},
				ConvertField[time.Time, string](
					"BirthDay",
					converter.ConverterFunc[time.Time, string](func(src time.Time) (string, error) {
						return val, nil
					})))
			assert.Nil(t, err)
			assert.Equal(t, &ConvSimpleDst{BirthDay: val}, c)
		}(i)
	}
	wg.Wait()
}

type BasicSrc struct {
	Name    string
	Age     int
	CNumber complex64
}

type BasicDst struct {
	Name    string
	Age     int
	CNumber complex64
}

type SimpleSrc struct {
	Name    string
	Age     *int
	Friends []string
}

type SimpleDst struct {
	Name    string
	Age     *int
	Friends []string
}

type EmbedSrc struct {
	SimpleSrc
	*BasicSrc
}

type EmbedDst struct {
	SimpleSrc
	*BasicSrc
}

type ComplexSrc struct {
	Simple SimpleSrc
	Embed  *EmbedSrc
	BasicSrc
}

type ComplexDst struct {
	Simple SimpleDst
	Embed  *EmbedDst
	BasicSrc
}

type SpecialSrc struct {
	Arr [3]float32
	M   map[string]int
}

type SpecialDst struct {
	Arr [3]float32
	M   map[string]int
}

type InterfaceSrc interface {
}

type InterfaceDst interface {
}

type NotMatchSrc struct {
	Simple SimpleSrc
	Embed  *EmbedSrc
	BasicSrc
	S struct {
		A string
	}
}

type NotMatchDst struct {
	Simple SimpleDst
	Embed  *EmbedDst
	BasicSrc
	S struct {
		A int
	}
}

type MultiPtrSrc struct {
	Name    string
	Age     **int
	Friends []string
}

type MultiPtrDst struct {
	Name    string
	Age     **int
	Friends []string
}

type DiffSrc struct {
	A string
	B int
	c SimpleSrc
	F BasicSrc
}
type DiffDst struct {
	A string
	B int
	d SimpleSrc
	G BasicSrc
}

type SimpleEmbedDst struct {
	SimpleSrc
}

type ArraySrc struct {
	A []SimpleSrc
}

type ArrayDst struct {
	A []SimpleSrc
}

type ArrayDst1 struct {
	A []SimpleDst
}

type MapSrc struct {
	A map[string]SimpleSrc
}

type MapDst struct {
	A map[string]SimpleSrc
}

type MapDst1 struct {
	A map[string]SimpleDst
}

type SpecialSrc1 struct {
	A int
}

type aliasInt int
type SpecialDst1 struct {
	A aliasInt
}

type aliasInt1 = int
type SpecialDst2 struct {
	A aliasInt1
}

type ConvSimpleSrc struct {
	Name     string
	Age      *int
	BirthDay time.Time
	Friends  []string
}

type ConvSimpleDst struct {
	Name     string
	Age      *int
	BirthDay string
	Friends  []string
}

type ConvSpecialSrc struct {
	Arr  [3]float32
	M    map[string]int
	Diff map[string]int
}

type ConvSpecialDst struct {
	Arr  [3]float32
	M    map[string]int
	Diff []int
}

func BenchmarkReflectCopier_Copy(b *testing.B) {
	// 复用 Copier
	b.Run("reused", func(b *testing.B) {
		copier, err := NewReflectCopier[SimpleSrc, SimpleDst]()
		if err != nil {
			b.Fatal(err)
		}
		for i := 1; i <= b.N; i++ {
			_, _ = copier.Copy(&SimpleSrc{
				Name:    "大明",
				Age:     ekit.ToPtr[int](18),
				Friends: []string{"Tom", "Jerry"},
			})
		}
	})

	// 每次都是新建
	b.Run("create", func(b *testing.B) {
		for i := 1; i <= b.N; i++ {
			copier, _ := NewReflectCopier[SimpleSrc, SimpleDst]()
			_, _ = copier.Copy(&SimpleSrc{
				Name:    "大明",
				Age:     ekit.ToPtr[int](18),
				Friends: []string{"Tom", "Jerry"},
			})
		}
	})
}

func BenchmarkReflectCopier_CopyComplexStruct(b *testing.B) {
	b.Run("reused", func(b *testing.B) {
		copier, _ := NewReflectCopier[ComplexSrc, ComplexDst]()
		for i := 1; i <= b.N; i++ {
			_, _ = copier.Copy(&ComplexSrc{
				Simple: SimpleSrc{
					Name:    "xiaohong",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"ha", "ha", "le"},
				},
				Embed: &EmbedSrc{
					SimpleSrc: SimpleSrc{
						Name:    "xiaopeng",
						Age:     ekit.ToPtr[int](88),
						Friends: []string{"la", "ha", "le"},
					},
					BasicSrc: &BasicSrc{
						Name:    "wang",
						Age:     22,
						CNumber: complex(2, 1),
					},
				},
				BasicSrc: BasicSrc{
					Name:    "wang11",
					Age:     22,
					CNumber: complex(2, 1),
				},
			})
		}
	})
	b.Run("create", func(b *testing.B) {
		for i := 1; i <= b.N; i++ {
			copier, _ := NewReflectCopier[ComplexSrc, ComplexDst]()
			_, _ = copier.Copy(&ComplexSrc{
				Simple: SimpleSrc{
					Name:    "xiaohong",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"ha", "ha", "le"},
				},
				Embed: &EmbedSrc{
					SimpleSrc: SimpleSrc{
						Name:    "xiaopeng",
						Age:     ekit.ToPtr[int](88),
						Friends: []string{"la", "ha", "le"},
					},
					BasicSrc: &BasicSrc{
						Name:    "wang",
						Age:     22,
						CNumber: complex(2, 1),
					},
				},
				BasicSrc: BasicSrc{
					Name:    "wang11",
					Age:     22,
					CNumber: complex(2, 1),
				},
			})
		}
	})
}
