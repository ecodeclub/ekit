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

package sqlx

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewNullBool(t *testing.T) {
	tests := []struct {
		name string
		val  bool
		want sql.NullBool
	}{
		{
			name: "nonzero",
			val:  true,
			want: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		{
			name: "zero",
			val:  false,
			want: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullBool(tt.val), "NewNullBool(%v)", tt.val)
		})
	}
}

func TestNewNullBytes(t *testing.T) {
	tests := []struct {
		name string
		val  []byte
		want sql.NullString
	}{
		{
			name: "nonzero",
			val:  []byte("test"),
			want: sql.NullString{
				String: "test",
				Valid:  true,
			},
		},
		{
			name: "zero",
			val:  []byte{},
			want: sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullBytes(tt.val), "NewNullBytes(%v)", tt.val)
		})
	}
}

func TestNewNullFloat64(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want sql.NullFloat64
	}{
		{
			name: "nonzero",
			val:  1.1,
			want: sql.NullFloat64{
				Float64: 1.1,
				Valid:   true,
			},
		},
		{
			name: "zero",
			val:  0,
			want: sql.NullFloat64{
				Float64: 0,
				Valid:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullFloat64(tt.val), "NewNullFloat64(%v)", tt.val)
		})
	}
}

func TestNewNullInt64(t *testing.T) {
	tests := []struct {
		name string
		val  int64
		want sql.NullInt64
	}{
		{
			name: "nonzero",
			val:  1,
			want: sql.NullInt64{
				Int64: 1,
				Valid: true,
			},
		},
		{
			name: "zero",
			val:  0,
			want: sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullInt64(tt.val), "NewNullInt64(%v)", tt.val)
		})
	}
}

func TestNewNullString(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want sql.NullString
	}{
		{
			name: "nonzero",
			val:  "test",
			want: sql.NullString{
				String: "test",
				Valid:  true,
			},
		},
		{
			name: "zero",
			val:  "",
			want: sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullString(tt.val), "NewNullString(%v)", tt.val)
		})
	}
}

func TestNewNullTime(t *testing.T) {
	tests := []struct {
		name string
		val  time.Time
		want sql.NullTime
	}{
		{
			name: "nonzero",
			val:  time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			want: sql.NullTime{
				Time:  time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
				Valid: true,
			},
		},
		{
			name: "zero",
			val:  time.Time{},
			want: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullTime(tt.val), "NewNullTime(%v)", tt.val)
		})
	}
}
