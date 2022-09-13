package copier

import (
	"github.com/gotomicro/ekit"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// 注意！！
// 测试用例是在arm64的环境下人为计算得到的
func TestStructHelper(t *testing.T) {
	testCases := []struct {
		name     string
		inStruct any
		wantDst  *structOffsets
		wantErr  error
	}{
		{
			name:     "错误输入",
			inStruct: int32(2),
			wantDst: &structOffsets{
				helper: []structHelper{},
			},
			wantErr: errorNotStruct,
		},
		{
			name:     "没有指针",
			inStruct: NonePtr{},
			wantDst: &structOffsets{
				helper: []structHelper{},
			},
			wantErr: nil,
		},
		{
			name:     "SimplePtr1",
			inStruct: SimplePtr1{},
			wantDst: &structOffsets{
				helper: []structHelper{
					{8, reflect.TypeOf(ekit.ToPtr[int64](1))},
					{16, reflect.TypeOf(ekit.ToPtr[int32](1))},
				},
			},
			wantErr: nil,
		},
		{
			name:     "SimplePtr2",
			inStruct: SimplePtr2{},
			wantDst: &structOffsets{
				helper: []structHelper{
					{8, reflect.TypeOf(ekit.ToPtr[int64](1))},
					{16, reflect.TypeOf(ekit.ToPtr[int32](1))},
				},
			},
			wantErr: nil,
		},
		{
			name:     "SimplePtr3",
			inStruct: SimplePtr3{},
			wantDst: &structOffsets{
				helper: []structHelper{
					{8, reflect.TypeOf(ekit.ToPtr[int64](1))},
					{16, reflect.TypeOf(ekit.ToPtr[int32](1))},
				},
			},
			wantErr: nil,
		},
		{
			name:     "Composite1",
			inStruct: Composite1{},
			wantDst: &structOffsets{
				helper: []structHelper{
					{16, reflect.TypeOf(ekit.ToPtr[int64](1))},
					{24, reflect.TypeOf(ekit.ToPtr[int32](1))},
				},
			},
			wantErr: nil,
		},
		{
			name:     "Composite2",
			inStruct: Composite2{},
			wantDst: &structOffsets{
				helper: []structHelper{
					{32, reflect.TypeOf(ekit.ToPtr[int64](1))},
					{40, reflect.TypeOf(ekit.ToPtr[int32](1))},
				},
			},
			wantErr: nil,
		},
		{
			name:     "Composite3",
			inStruct: Composite3{},
			wantDst: &structOffsets{
				helper: []structHelper{
					{40, reflect.TypeOf(ekit.ToPtr[int64](1))},
					{48, reflect.TypeOf(ekit.ToPtr[int32](1))},
				},
			},
			wantErr: nil,
		},
		{
			name:     "SpecialOffsets",
			inStruct: SpecialOffsets{},
			wantDst: &structOffsets{
				helper: []structHelper{},
				deepCopyOffsets: map[reflect.Kind][]uintptr{
					/*
						type SliceHeader struct {
							Data uintptr
							Len  int
							Cap  int
						}
					*/
					reflect.Slice: []uintptr{0},
					reflect.Map:   []uintptr{24},
					reflect.Array: []uintptr{32},
				},
			},
			wantErr: nil,
		},
	}

	helperMap := make(map[string]*structOffsets)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := FindOffsets(tc.inStruct, helperMap)
			if res != nil {
				assert.Equal(t, tc.wantDst, res)
			}
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

type NonePtr struct {
	b int64
	a int32
	c uint32
}

type SimplePtr1 struct {
	a int64
	b *int64
	c *int32
}

type SimplePtr2 struct {
	a  int32
	a2 int32
	b  *int64
	c  *int32
}

// golang的对齐操作
type SimplePtr3 struct {
	a int32
	b *int64
	c *int32
}

type Composite1 struct {
	a      int32
	simple SimplePtr1
}

type Composite2 struct {
	NonePtr
	Composite1
}

type Composite3 struct {
	a int32
	b int64
	c uint32
	Composite1
}

type SpecialOffsets struct {
	x []int
	y map[string][]int
	z [5]int
}
