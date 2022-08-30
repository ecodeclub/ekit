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
	"reflect"
	"testing"

	"github.com/gotomicro/ekit"
	"github.com/stretchr/testify/assert"
)

func TestReflectCopier_Copy(t *testing.T) {
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
			wantErr: newErrKindNotMatchError(reflect.String, reflect.Int, "A"),
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
			wantDst: &ArrayDst1{},
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

func BenchmarkReflectCopier_Copy(b *testing.B) {
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
}

func BenchmarkReflectCopier_CopyComplexStruct(b *testing.B) {
	copier, err := NewReflectCopier[ComplexSrc, ComplexDst]()
	if err != nil {
		b.Fatal(err)
	}
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
}
