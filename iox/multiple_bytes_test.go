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

package iox

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleBytesReadWrite(t *testing.T) {
	testCases := []struct {
		name     string
		write    []byte
		readSize int
		wantRead []byte
		wantErr  error
	}{
		{
			name:     "empty read",
			write:    []byte{},
			readSize: 1,
			wantRead: []byte{},
			wantErr:  io.EOF,
		},
		{
			name:     "single byte",
			write:    []byte{1},
			readSize: 1,
			wantRead: []byte{1},
			wantErr:  nil,
		},
		{
			name:     "multiple bytes",
			write:    []byte{1, 2, 3, 4, 5},
			readSize: 3,
			wantRead: []byte{1, 2, 3},
			wantErr:  nil,
		},
		{
			name:     "read more than available",
			write:    []byte{1, 2},
			readSize: 4,
			wantRead: []byte{1, 2},
			wantErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mb := NewMultipleBytes(len(tc.write))
			n, err := mb.Write(tc.write)
			assert.Equal(t, len(tc.write), n)
			assert.Nil(t, err)

			read := make([]byte, tc.readSize)
			n, err = mb.Read(read)
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.wantRead, read[:n])
			}
		})
	}
}

func TestMultipleBytesReadEdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		writes    [][]byte // 多次写入的数据
		readSizes []int    // 多次读取的大小
		wantReads [][]byte // 期望的读取结果
		wantErrs  []error  // 期望的错误
	}{
		{
			name:      "单片-恰好读完",
			writes:    [][]byte{{1, 2, 3}},
			readSizes: []int{3},
			wantReads: [][]byte{{1, 2, 3}},
			wantErrs:  []error{nil},
		},
		{
			name:      "单片-部分读取",
			writes:    [][]byte{{1, 2, 3}},
			readSizes: []int{2, 1},
			wantReads: [][]byte{{1, 2}, {3}},
			wantErrs:  []error{nil, nil},
		},
		{
			name:      "单片-读取溢出",
			writes:    [][]byte{{1, 2}},
			readSizes: []int{3},
			wantReads: [][]byte{{1, 2}},
			wantErrs:  []error{nil},
		},
		{
			name:      "多片-跨片读取",
			writes:    [][]byte{{1, 2}, {3, 4}, {5, 6}},
			readSizes: []int{4},
			wantReads: [][]byte{{1, 2, 3, 4}},
			wantErrs:  []error{nil},
		},
		{
			name:      "多片-恰好读完",
			writes:    [][]byte{{1, 2}, {3, 4}},
			readSizes: []int{4},
			wantReads: [][]byte{{1, 2, 3, 4}},
			wantErrs:  []error{nil},
		},
		{
			name:      "多片-未读完",
			writes:    [][]byte{{1, 2}, {3, 4}, {5, 6}},
			readSizes: []int{3},
			wantReads: [][]byte{{1, 2, 3}},
			wantErrs:  []error{nil},
		},
		{
			name:      "索引边界-首尾交叉验证",
			writes:    [][]byte{{1}, {2}, {3}},
			readSizes: []int{1, 1, 1, 1},
			wantReads: [][]byte{{1}, {2}, {3}, {}},
			wantErrs:  []error{nil, nil, nil, io.EOF},
		},
		{
			name:      "空切片读取",
			writes:    [][]byte{{}},
			readSizes: []int{1},
			wantReads: [][]byte{{}},
			wantErrs:  []error{io.EOF},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mb := NewMultipleBytes(len(tc.writes))

			// 写入数据
			for _, w := range tc.writes {
				n, err := mb.Write(w)
				assert.Equal(t, len(w), n)
				assert.Nil(t, err)
			}

			// 读取数据
			for i, size := range tc.readSizes {
				read := make([]byte, size)
				n, err := mb.Read(read)
				if tc.wantErrs[i] != nil {
					assert.Equal(t, tc.wantErrs[i], err)
				} else {
					assert.Nil(t, err)
					assert.Equal(t, tc.wantReads[i], read[:n])
				}
			}
		})
	}
}

func TestMultipleBytesReset(t *testing.T) {
	mb := NewMultipleBytes(4)
	data := []byte{1, 2, 3, 4}

	// 写入数据
	n, err := mb.Write(data)
	assert.Equal(t, len(data), n)
	assert.Nil(t, err)

	// 第一次读取
	read := make([]byte, 2)
	n, err = mb.Read(read)
	assert.Equal(t, 2, n)
	assert.Nil(t, err)
	assert.Equal(t, []byte{1, 2}, read)

	// 重置
	mb.Reset()

	// 重置后再次读取
	read = make([]byte, 4)
	n, err = mb.Read(read)
	assert.Equal(t, 4, n)
	assert.Nil(t, err)
	assert.Equal(t, data, read)
}
