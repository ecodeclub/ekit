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
	"reflect"
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

func TestReflectCopier_CopyTo(t *testing.T) {
	testCases := []struct {
		name     string
		copyFunc func() (any, error)
		wantDst  any
		wantErr  error
	}{
		{
			name: "simple struct",
			copyFunc: func() (any, error) {
				dst := &SimpleDst{}
				err := CopyTo(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				}, dst)
				return dst, err
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
				dst := &BasicDst{}
				err := CopyTo(&BasicSrc{
					Name:    "大明",
					Age:     10,
					CNumber: complex(1, 2),
				}, dst)
				return dst, err
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
				i := 10
				dst := ekit.ToPtr(int(0))
				err := CopyTo(&i, dst)
				return dst, err
			},
			wantErr: newErrTypeError(reflect.TypeOf(10)),
		},
		{
			name: "dst 是基础类型",
			copyFunc: func() (any, error) {

				dst := ekit.ToPtr("")
				err := CopyTo(&SimpleSrc{
					Name:    "大明",
					Age:     ekit.ToPtr[int](18),
					Friends: []string{"Tom", "Jerry"},
				}, dst)
				return dst, err
			},
			wantErr: newErrTypeError(reflect.TypeOf("")),
		},
		{
			name: "接口类型",
			copyFunc: func() (any, error) {
				i := InterfaceSrc(10)
				dst := ekit.ToPtr(InterfaceDst(10))
				err := CopyTo(&i, dst)
				return dst, err
			},
			wantErr: newErrTypeError(reflect.TypeOf(new(InterfaceSrc)).Elem()),
		},
		{
			name: "simple struct 空切片, 空指针",
			copyFunc: func() (any, error) {
				dst := &SimpleDst{}
				err := CopyTo(&SimpleSrc{
					Name: "大明",
				}, dst)
				return dst, err
			},
			wantDst: &SimpleDst{
				Name: "大明",
			},
		},
		{
			name: "组合 struct ",
			copyFunc: func() (any, error) {
				dst := &EmbedDst{}
				err := CopyTo(&EmbedSrc{
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
				}, dst)
				return dst, err
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
				dst := &ComplexDst{}
				err := CopyTo(&ComplexSrc{
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
				}, dst)
				return dst, err
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
				dst := &SpecialDst{}
				err := CopyTo(&SpecialSrc{
					Arr: [3]float32{1, 2, 3},
					M: map[string]int{
						"ha": 1,
						"o":  2,
					},
				}, dst)
				return dst, err
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
				dst := &NotMatchDst{}
				err := CopyTo(&NotMatchSrc{
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
				}, dst)
				return dst, err
			},
			wantErr: newErrKindNotMatchError(reflect.String, reflect.Int, "A"),
		},
		{
			name: "多重指针",
			copyFunc: func() (any, error) {
				dst := &MultiPtrDst{}
				err := CopyTo(&MultiPtrSrc{
					Name:    "a",
					Age:     ekit.ToPtr[*int](ekit.ToPtr[int](10)),
					Friends: nil,
				}, dst)
				return dst, err
			},
			wantErr: newErrMultiPointer("Age"),
		},
		{
			name: "src 有额外字段",
			copyFunc: func() (any, error) {
				dst := &DiffDst{}
				err := CopyTo(&DiffSrc{
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
				d: SimpleSrc{},
				G: BasicSrc{},
			},
		},
		{
			name: "dst 有额外字段",
			copyFunc: func() (any, error) {

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
				err := CopyTo(&DiffSrc{
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
				dst := &SimpleEmbedDst{}
				err := CopyTo(&SimpleSrc{
					Name: "haha",
				}, dst)
				return dst, err
			},
			wantDst: &SimpleEmbedDst{},
		},
		{
			name: "成员为结构体数组",
			copyFunc: func() (any, error) {
				dst := &ArrayDst{}
				return dst, CopyTo(&ArraySrc{
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
				}, dst)
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
				dst := &ArrayDst1{}
				return dst, CopyTo(&ArraySrc{
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
				}, dst)
			},
			wantErr: newErrTypeNotMatchError(reflect.TypeOf(new([]SimpleSrc)).Elem(), reflect.TypeOf(new([]SimpleDst)).Elem(), "A"),
		},
		{
			name: "成员为map结构体",
			copyFunc: func() (any, error) {
				dst := &MapDst{}
				return dst, CopyTo(&MapSrc{
					A: map[string]SimpleSrc{
						"a": {
							Name:    "大明",
							Age:     ekit.ToPtr[int](18),
							Friends: []string{"Tom", "Jerry"},
						},
					},
				}, dst)
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
				dst := &MapDst1{}
				return dst, CopyTo(&MapSrc{
					A: map[string]SimpleSrc{
						"a": {
							Name:    "大明",
							Age:     ekit.ToPtr[int](18),
							Friends: []string{"Tom", "Jerry"},
						},
					},
				}, dst)
			},
			wantErr: newErrTypeNotMatchError(reflect.TypeOf(new(map[string]SimpleSrc)).Elem(), reflect.TypeOf(new(map[string]SimpleDst)).Elem(), "A"),
		},
		{
			name: "成员有别名类型",
			copyFunc: func() (any, error) {
				dst := &SpecialDst1{}
				return dst, CopyTo(&SpecialSrc1{
					A: 1,
				}, dst)
			},
			wantErr: newErrTypeNotMatchError(reflect.TypeOf(new(int)).Elem(), reflect.TypeOf(new(aliasInt)).Elem(), "A"),
		},
		{
			name: "成员有别名类型1",
			copyFunc: func() (any, error) {
				dst := &SpecialDst2{}
				return dst, CopyTo(&SpecialSrc1{
					A: 1,
				}, dst)
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

func BenchmarkReflectCopier_Copy_PureRunTime(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		_ = CopyTo(&SimpleSrc{
			Name:    "大明",
			Age:     ekit.ToPtr[int](18),
			Friends: []string{"Tom", "Jerry"},
		}, &SimpleDst{})
	}
}

func BenchmarkReflectCopier_CopyComplexStruct_WithPureRuntime(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		_ = CopyTo(&ComplexSrc{
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
		}, &ComplexDst{})
	}
}
