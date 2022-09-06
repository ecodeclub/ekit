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

package copier

import (
	"github.com/gotomicro/ekit"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"unsafe"
)

func TestReflectCopier_Copy(t *testing.T) {
	testCases := []struct {
		name     string
		copyFunc func() (any, error)
		wantDst  any
		wantErr  error
	}{
		{
			name: "error input",
			copyFunc: func() (any, error) {
				copier, _ := NewReflectCopier[*SimpleSrc, SimpleDst]()
				simpleSrc := &SimpleSrc{}
				return copier.Copy(&simpleSrc)
			},
			wantErr: errInvalidType,
		},
		{
			name: "error input int",
			copyFunc: func() (any, error) {
				copier, _ := NewReflectCopier[int, SimpleDst]()
				test := 1
				return copier.Copy(&test)
			},
			wantErr: errInvalidType,
		},
		// Map、Slice的复杂结构应不予以支持深度拷贝
		// 这种"二级"且内存空间较大的引用(本质还是指针)数据结构，
		// 从内存、性能、安全上来看，需要尽可能地减少深度拷贝
		// 所以应当把Map和Slice的管理交由调用者,从实际业务角度出发，目前不太清楚深拷贝一个Map和slice干什么，真要有，不如让具体业务自己深度拷贝
		// 所以选择直接浅拷贝
		// 而String类型，约等于C++中的 string_view，所以直接Set就可以，即dst中的string获得了底层字符串的一份视图
		// http://t.csdn.cn/iMzqZ
		{
			name: "simple struct",
			copyFunc: func() (any, error) {
				copier, _ := NewReflectCopier[SimpleInlineSrc, SimpleInlineDst]()
				initSimpleStructParam()
				res, err := copier.Copy(simpleInlineSrc)
				//修改simpleInlineSrc中的string,(更换原先src中的string，即src中的string_view是 "abc"，具体见函数)
				changeSimpleStructParam()
				return res, err
			},
			wantDst: &SimpleInlineDst{
				mapS:         map[string]int{"1": 1, "2": 2, "3": 3, "4": 4},
				byteS:        []byte{'a', 'b', 'c', 'd'},
				stringS:      "before",
				stringIgnore: "",
			},
		},
		{
			name: "simple struct demo",
			copyFunc: func() (any, error) {
				copier, _ := NewReflectCopier[SimpleSrc, SimpleDst]()
				age := 18
				simpleSrc := &SimpleSrc{
					Name:    "大明大聪明",
					Age:     &age,
					Friends: []string{"Tom", "Jerry"},
				}
				res, err := copier.Copy(simpleSrc)
				//都进行修改再比较，以验证深度拷贝
				*simpleSrc.Age = 19
				simpleSrc.Name = "大明大xx"
				assert.Equal(t, 19, age)
				return res, err
			},
			wantDst: &SimpleDst{
				Name:    "大明大聪明",
				Age:     ekit.ToPtr[int](18),
				Friends: []string{"Tom", "Jerry"},
			},
		},
		{
			name: "pointer struct",
			copyFunc: func() (any, error) {
				copier, _ := NewReflectCopier[SimplePointersSrc, SimplePointersDst]()
				a := 1
				b := float32(1.0)
				c := float64(2.0)
				d := uint8(1)
				e := uint32(2)
				f := uint64(3)
				g := int64(-1)
				h := int32(-2)
				i := int8(-3)
				j := bool(false)
				k := complex64(complex(1, 2))
				l := complex(3, 4)

				simpleSrc := &SimplePointersSrc{
					&a, &b, &c, &d, &e, &f, &g, &h, &i, &j, &k, &l,
				}
				res, err := copier.Copy(simpleSrc)

				a = 2
				b = float32(2.0)
				c = float64(3.0)
				d = uint8(2)
				e = uint32(4)
				f = uint64(5)
				g = int64(-2)
				h = int32(-3)
				i = int8(-4)
				j = bool(true)
				k = complex64(complex(2, 3))
				l = complex(4, 5)

				return res, err
			},
			wantDst: &SimplePointersDst{
				ekit.ToPtr[int](1),
				ekit.ToPtr[float32](1.0),
				ekit.ToPtr[float64](2.0),
				ekit.ToPtr[uint8](1),
				ekit.ToPtr[uint32](2),
				ekit.ToPtr[uint64](3),
				ekit.ToPtr[int64](-1),
				ekit.ToPtr[int32](-2),
				ekit.ToPtr[int8](-3),
				ekit.ToPtr[bool](false),
				ekit.ToPtr[complex64](complex64(complex(1, 2))),
				ekit.ToPtr[complex128](complex(3, 4)),
				nil,
			},
		}, {
			name: "组合,子数据结构内部全为指针",
			copyFunc: func() (any, error) {
				copier, _ := NewReflectCopier[CompositeSrc1, CompositeDst1]()
				a := 1
				b := float32(1.0)
				c := float64(2.0)
				d := uint8(1)
				e := uint32(2)
				f := uint64(3)
				g := int64(-1)
				h := int32(-2)
				i := int8(-3)
				j := bool(false)
				k := complex64(complex(1, 2))
				l := complex(3, 4)
				simpleSrc := SimplePointersSrc{
					&a, &b, &c, &d, &e, &f, &g, &h, &i, &j, &k, &l,
				}
				src1 := &CompositeSrc1{
					Simple: simpleSrc,
					a:      1,
					b:      1,
				}
				res, err := copier.Copy(src1)

				a = 2
				b = float32(2.0)
				c = float64(3.0)
				d = uint8(2)
				e = uint32(4)
				f = uint64(5)
				g = int64(-2)
				h = int32(-3)
				i = int8(-4)
				j = bool(true)
				k = complex64(complex(2, 3))
				l = complex(4, 5)

				//修改simpleInlineSrc中的string,(更换原先src中的string，即src中的string_view是 "abc"，具体见函数)
				return res, err
			},
			wantDst: &CompositeDst1{
				a: 1,
				Simple: SimplePointersSrc{
					ekit.ToPtr[int](1),
					ekit.ToPtr[float32](1.0),
					ekit.ToPtr[float64](2.0),
					ekit.ToPtr[uint8](1),
					ekit.ToPtr[uint32](2),
					ekit.ToPtr[uint64](3),
					ekit.ToPtr[int64](-1),
					ekit.ToPtr[int32](-2),
					ekit.ToPtr[int8](-3),
					ekit.ToPtr[bool](false),
					ekit.ToPtr[complex64](complex64(complex(1, 2))),
					ekit.ToPtr[complex128](complex(3, 4)),
				},
			},
		}, {
			name: "Composite2",
			copyFunc: func() (any, error) {
				copier, _ := NewReflectCopier[CompositeSrc2, CompositeDst2]()
				a := 1
				b := float32(1.0)
				c := float64(2.0)
				d := uint8(1)
				e := uint32(2)
				f := uint64(3)
				g := int64(-1)
				h := int32(-2)
				i := int8(-3)
				j := bool(false)
				k := complex64(complex(1, 2))
				l := complex(3, 4)
				simpleSrc := &SimplePointersSrc{
					&a, &b, &c, &d, &e, &f, &g, &h, &i, &j, &k, &l,
				}
				src1 := &CompositeSrc2{
					Simple: simpleSrc,
					a:      1,
					b:      1,
				}
				res, err := copier.Copy(src1)

				a = 2
				b = float32(2.0)
				c = float64(3.0)
				d = uint8(2)
				e = uint32(4)
				f = uint64(5)
				g = int64(-2)
				h = int32(-3)
				i = int8(-4)
				j = bool(true)
				k = complex64(complex(2, 3))
				l = complex(4, 5)
				//修改simpleInlineSrc中的string,(更换原先src中的string，即src中的string_view是 "abc"，具体见函数)
				return res, err
			},
			wantDst: &CompositeDst2{
				a: 1,
				Simple: &SimplePointersSrc{
					ekit.ToPtr[int](1),
					ekit.ToPtr[float32](1.0),
					ekit.ToPtr[float64](2.0),
					ekit.ToPtr[uint8](1),
					ekit.ToPtr[uint32](2),
					ekit.ToPtr[uint64](3),
					ekit.ToPtr[int64](-1),
					ekit.ToPtr[int32](-2),
					ekit.ToPtr[int8](-3),
					ekit.ToPtr[bool](false),
					ekit.ToPtr[complex64](complex64(complex(1, 2))),
					ekit.ToPtr[complex128](complex(3, 4)),
				},
			},
		},
		// 你还需要测试
		// 1. Src 或者 Dst 类型非法，例如基本类型，内置类型或者接口
		// 2. 测试组合（结构体组合，指针组合，接口组合——接口组合可以直接不支持），深层组合，多重组合
		// 3. 复杂类型字段，如字段是结构体，字段是结构体指针，以及多级指针（不需要支持）
		// 4. 类型不匹配
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

type SimpleInlineSrc struct {
	mapS    map[string]int
	byteS   []byte
	stringS string
}

type SimpleInlineDst struct {
	mapS         map[string]int
	byteS        []byte
	stringS      string
	stringIgnore string
}

var simpleMap = make(map[string]int)
var simpleSlice = make([]byte, 4)
var simpleInlineSrc = &SimpleInlineSrc{
	mapS:  simpleMap,
	byteS: simpleSlice}

func initSimpleStructParam() {
	simpleMap["1"] = 1
	simpleMap["2"] = 2
	simpleMap["3"] = 3
	simpleSlice[0] = 'a'
	simpleSlice[1] = 'b'
	simpleSlice[2] = 'c'
	simpleInlineSrc.stringS = "before"
}

func changeSimpleStructParam() {
	simpleMap["4"] = 4
	simpleSlice[3] = 'd'
	simpleInlineSrc.stringS = "after"
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

type SimplePointersSrc struct {
	a *int
	b *float32
	c *float64
	d *uint8
	e *uint32
	f *uint64
	g *int64
	h *int32
	i *int8
	j *bool
	k *complex64
	l *complex128
}

type SimplePointersDst struct {
	a *int
	b *float32
	c *float64
	d *uint8
	e *uint32
	f *uint64
	g *int64
	h *int32
	i *int8
	j *bool
	k *complex64
	l *complex128
	m *unsafe.Pointer
}

type CompositeSrc1 struct {
	Simple SimplePointersSrc
	a      int
	b      complex64
}

type CompositeDst1 struct {
	Simple SimplePointersSrc
	a      int
}

type CompositeSrc2 struct {
	Simple *SimplePointersSrc
	a      int
	b      complex64
}

type CompositeDst2 struct {
	Simple *SimplePointersSrc
	a      int
}

func TestReflectCopier_LongYue(t *testing.T) {
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
				copier, err := NewReflectCopier[SimpleSrcLongYue, SimpleDstLongYue]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&SimpleSrcLongYue{
					Name: "大明",
				})
			},
			wantDst: &SimpleDstLongYue{
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
					SimpleSrcLongYue: SimpleSrcLongYue{
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
				SimpleSrcLongYue: SimpleSrcLongYue{
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
		//只支持名称和类型完全相同的字段
		{
			name: "复杂 Struct",
			copyFunc: func() (any, error) {
				copier, err := NewReflectCopier[ComplexSrc, ComplexDst]()
				if err != nil {
					return nil, err
				}
				return copier.Copy(&ComplexSrc{
					Simple: SimpleSrcLongYue{
						Name:    "xiaohong",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrcLongYue: SimpleSrcLongYue{
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
				Embed: &EmbedSrc{
					SimpleSrcLongYue: SimpleSrcLongYue{
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
					Simple: SimpleSrcLongYue{
						Name:    "xiaohong",
						Age:     ekit.ToPtr[int](18),
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrcLongYue: SimpleSrcLongYue{
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
			wantDst: &NotMatchDst{
				Embed: &EmbedSrc{
					SimpleSrcLongYue: SimpleSrcLongYue{
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
			wantErr: nil,
		},
		{
			name: "支持多重指针",
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
			wantDst: &MultiPtrDst{
				Name:    "a",
				Age:     ekit.ToPtr[*int](ekit.ToPtr[int](10)),
				Friends: nil,
			},
			wantErr: nil,
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
			name: "成员为结构体数组，目前仅为浅拷贝，其实这儿的测试用例没啥意义",
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
			name: "成员为结构体数组，结构体不同，不返回错误，直接忽略该字段",
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
			wantDst: &ArrayDst1{},
			wantErr: nil,
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
			wantDst: &MapDst1{},
			wantErr: nil,
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
			wantDst: &SpecialDst1{},
			wantErr: nil,
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

type SimpleSrcLongYue struct {
	Name    string
	Age     *int
	Friends []string
}

type SimpleDstLongYue struct {
	Name    string
	Age     *int
	Friends []string
}

type EmbedSrc struct {
	SimpleSrcLongYue
	*BasicSrc
}

type EmbedDst struct {
	SimpleSrcLongYue
	*BasicSrc
}

type ComplexSrc struct {
	Simple SimpleSrcLongYue
	Embed  *EmbedSrc
	BasicSrc
}

type ComplexDst struct {
	Simple SimpleDstLongYue
	Embed  *EmbedSrc
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
	Simple SimpleSrcLongYue
	Embed  *EmbedSrc
	BasicSrc
	S struct {
		A string
	}
}

type NotMatchDst struct {
	Simple SimpleDstLongYue
	Embed  *EmbedSrc
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
