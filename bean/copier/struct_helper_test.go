package copier

import (
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
				ptrOffsets: []uintptr{},
			},
			wantErr: errorNotStruct,
		},
		{
			name:     "没有指针",
			inStruct: NonePtr{},
			wantDst: &structOffsets{
				ptrOffsets: []uintptr{},
			},
			wantErr: nil,
		},
		{
			name:     "SimplePtr1",
			inStruct: SimplePtr1{},
			wantDst: &structOffsets{
				ptrOffsets: []uintptr{8, 16},
			},
			wantErr: nil,
		},
		{
			name:     "SimplePtr2",
			inStruct: SimplePtr2{},
			wantDst: &structOffsets{
				ptrOffsets: []uintptr{8, 16},
			},
			wantErr: nil,
		},
		{
			name:     "SimplePtr3",
			inStruct: SimplePtr3{},
			wantDst: &structOffsets{
				ptrOffsets: []uintptr{8, 16},
			},
			wantErr: nil,
		},
		{
			name:     "Composite1",
			inStruct: Composite1{},
			wantDst: &structOffsets{
				ptrOffsets: []uintptr{16, 24},
			},
			wantErr: nil,
		},
		{
			name:     "Composite2",
			inStruct: Composite2{},
			wantDst: &structOffsets{
				ptrOffsets: []uintptr{32, 40},
			},
			wantErr: nil,
		},
		{
			name:     "Composite3",
			inStruct: Composite3{},
			wantDst: &structOffsets{
				ptrOffsets: []uintptr{40, 48},
			},
			wantErr: nil,
		},
		{
			name:     "SpecialOffsets",
			inStruct: SpecialOffsets{},
			wantDst: &structOffsets{
				ptrOffsets: []uintptr{},

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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := FindOffsetsDefault(tc.inStruct)
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
