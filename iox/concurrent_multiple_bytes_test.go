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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentMultipleBytes_ConcurrentReadWrite(t *testing.T) {
	testCases := []struct {
		name      string
		writes    [][]byte
		readSizes []int
		wantTotal int
	}{
		{
			name: "单片-恰好读完",
			writes: [][]byte{
				{1, 2, 3},
			},
			readSizes: []int{3},
			wantTotal: 3,
		},
		{
			name: "单片-未读完",
			writes: [][]byte{
				{1, 2, 3},
			},
			readSizes: []int{2},
			wantTotal: 2,
		},
		{
			name: "多片-跨片读取",
			writes: [][]byte{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			readSizes: []int{4},
			wantTotal: 4,
		},
		{
			name: "多片-恰好读完",
			writes: [][]byte{
				{1, 2},
				{3, 4},
			},
			readSizes: []int{4},
			wantTotal: 4,
		},
		{
			name: "多片-未读完",
			writes: [][]byte{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			readSizes: []int{3},
			wantTotal: 3,
		},
		{
			name: "索引边界-首尾交叉验证",
			writes: [][]byte{
				{1},
				{2},
				{3},
			},
			readSizes: []int{1, 1, 1},
			wantTotal: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmb := NewConcurrentMultipleBytes(len(tc.writes))

			// 并发写入
			wg := sync.WaitGroup{}
			for i := 0; i < len(tc.writes); i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					n, err := cmb.Write(tc.writes[idx])
					assert.Equal(t, len(tc.writes[idx]), n)
					assert.Nil(t, err)
				}(i)
			}
			wg.Wait()

			// 并发读取
			wg = sync.WaitGroup{}
			results := make([][]byte, len(tc.readSizes))
			for i := 0; i < len(tc.readSizes); i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					read := make([]byte, tc.readSizes[idx])
					n, err := cmb.Read(read)
					if err != nil && err != io.EOF {
						assert.Nil(t, err)
						return
					}
					results[idx] = read[:n]
				}(i)
			}
			wg.Wait()

			// 验证读取结果
			actualTotal := 0
			for _, res := range results {
				actualTotal += len(res)
			}
			assert.Equal(t, tc.wantTotal, actualTotal)
		})
	}
}

func TestConcurrentMultipleBytes_ConcurrentReset(t *testing.T) {
	cmb := NewConcurrentMultipleBytes(2)
	data := []byte{1, 2, 3, 4}

	// 写入初始数据
	n, err := cmb.Write(data)
	assert.Equal(t, len(data), n)
	assert.Nil(t, err)

	// 并发读取和重置
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(2)

		// 重置操作
		go func() {
			defer wg.Done()
			cmb.Reset()
		}()

		// 读取操作
		go func() {
			defer wg.Done()
			read := make([]byte, 4)
			n, err := cmb.Read(read)
			if err != nil && err != io.EOF {
				assert.Nil(t, err)
			}
			if n > 0 {
				assert.Equal(t, data[:n], read[:n])
			}
		}()
	}
	wg.Wait()
}

func TestConcurrentMultipleBytes_EdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		writes    [][]byte
		readSizes []int
		wantReads [][]byte
		wantErrs  []error
	}{
		{
			name:      "并发-空切片读取",
			writes:    [][]byte{{}},
			readSizes: []int{1},
			wantReads: [][]byte{{}},
			wantErrs:  []error{io.EOF},
		},
		{
			name:      "并发-多片读取",
			writes:    [][]byte{{1, 2}, {3, 4}, {5, 6}},
			readSizes: []int{2, 2, 2},
			wantReads: [][]byte{{1, 2}, {3, 4}, {5, 6}},
			wantErrs:  []error{nil, nil, nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmb := NewConcurrentMultipleBytes(len(tc.writes))

			// 并发写入
			wg := sync.WaitGroup{}
			for i := range tc.writes {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					n, err := cmb.Write(tc.writes[idx])
					assert.Equal(t, len(tc.writes[idx]), n)
					assert.Nil(t, err)
				}(i)
			}
			wg.Wait()

			// 并发读取
			results := make([][]byte, len(tc.readSizes))
			errs := make([]error, len(tc.readSizes))
			wg = sync.WaitGroup{}
			for i := range tc.readSizes {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					read := make([]byte, tc.readSizes[idx])
					n, err := cmb.Read(read)
					errs[idx] = err
					if err == nil {
						results[idx] = read[:n]
					}
				}(i)
			}
			wg.Wait()

			// 验证结果
			for i := range tc.wantReads {
				if tc.wantErrs[i] != nil {
					assert.Equal(t, tc.wantErrs[i], errs[i])
				} else {
					assert.Nil(t, errs[i])
					assert.Equal(t, tc.wantReads[i], results[i])
				}
			}
		})
	}
}
