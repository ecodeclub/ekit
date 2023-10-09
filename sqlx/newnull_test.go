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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewNullBool(t *testing.T) {
	type args struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want sql.NullBool
	}{
		{
			name: "test",
			args: args{
				val: true,
			},
			want: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		{
			name: "test",
			args: args{
				val: false,
			},
			want: sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullBool(tt.args.val), "NewNullBool(%v)", tt.args.val)
		})
	}
}

func TestNewNullBoolPtr(t *testing.T) {
	type args struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want *sql.NullBool
	}{
		{
			name: "test",
			args: args{
				val: true,
			},
			want: &sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		},
		{
			name: "test",
			args: args{
				val: false,
			},
			want: &sql.NullBool{
				Bool:  false,
				Valid: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullBoolPtr(tt.args.val), "NewNullBoolPtr(%v)", tt.args.val)
		})
	}
}

func TestNewNullBytes(t *testing.T) {
	type args struct {
		val []byte
	}
	tests := []struct {
		name string
		args args
		want sql.NullString
	}{
		{
			name: "test",
			args: args{
				val: []byte("test"),
			},
			want: sql.NullString{
				String: "test",
				Valid:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullBytes(tt.args.val), "NewNullBytes(%v)", tt.args.val)
		})
	}
}

func TestNewNullBytesPtr(t *testing.T) {
	type args struct {
		val []byte
	}
	tests := []struct {
		name string
		args args
		want *sql.NullString
	}{
		{
			name: "test",
			args: args{
				val: []byte("test"),
			},
			want: &sql.NullString{
				String: "test",
				Valid:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullBytesPtr(tt.args.val), "NewNullBytesPtr(%v)", tt.args.val)
		})
	}
}

func TestNewNullFloat64(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		args args
		want sql.NullFloat64
	}{
		{
			name: "test",
			args: args{
				val: 1.1,
			},
			want: sql.NullFloat64{
				Float64: 1.1,
				Valid:   true,
			},
		},
		{
			name: "test",
			args: args{
				val: 0,
			},
			want: sql.NullFloat64{
				Float64: 0,
				Valid:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullFloat64(tt.args.val), "NewNullFloat64(%v)", tt.args.val)
		})
	}
}

func TestNewNullFloat64Ptr(t *testing.T) {
	type args struct {
		val float64
	}
	tests := []struct {
		name string
		args args
		want *sql.NullFloat64
	}{
		{
			name: "test",
			args: args{
				val: 1.1,
			},
			want: &sql.NullFloat64{
				Float64: 1.1,
				Valid:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullFloat64Ptr(tt.args.val), "NewNullFloat64Ptr(%v)", tt.args.val)
		})
	}
}

func TestNewNullInt64(t *testing.T) {
	type args struct {
		val int64
	}
	tests := []struct {
		name string
		args args
		want sql.NullInt64
	}{
		{
			name: "test",
			args: args{
				val: 1,
			},
			want: sql.NullInt64{
				Int64: 1,
				Valid: true,
			},
		},
		{
			name: "test",
			args: args{
				val: 0,
			},
			want: sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullInt64(tt.args.val), "NewNullInt64(%v)", tt.args.val)
		})
	}
}

func TestNewNullInt64Ptr(t *testing.T) {
	type args struct {
		val int64
	}
	tests := []struct {
		name string
		args args
		want *sql.NullInt64
	}{
		{
			name: "test",
			args: args{
				val: 1,
			},
			want: &sql.NullInt64{
				Int64: 1,
				Valid: true,
			},
		},
		{
			name: "test",
			args: args{
				val: 0,
			},
			want: &sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullInt64Ptr(tt.args.val), "NewNullInt64Ptr(%v)", tt.args.val)
		})
	}
}

func TestNewNullString(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want sql.NullString
	}{
		{
			name: "test",
			args: args{
				val: "test",
			},
			want: sql.NullString{
				String: "test",
				Valid:  true,
			},
		},
		{
			name: "test",
			args: args{
				val: "",
			},
			want: sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullString(tt.args.val), "NewNullString(%v)", tt.args.val)
		})
	}
}

func TestNewNullStringPtr(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want *sql.NullString
	}{
		{
			name: "test",
			args: args{
				val: "test",
			},
			want: &sql.NullString{
				String: "test",
				Valid:  true,
			},
		},
		{
			name: "test",
			args: args{
				val: "",
			},
			want: &sql.NullString{
				String: "",
				Valid:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullStringPtr(tt.args.val), "NewNullStringPtr(%v)", tt.args.val)
		})
	}
}

func TestNewNullTime(t *testing.T) {
	type args struct {
		val time.Time
	}
	tests := []struct {
		name string
		args args
		want sql.NullTime
	}{
		{
			name: "test",
			args: args{
				val: time.Time{},
			},
			want: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullTime(tt.args.val), "NewNullTime(%v)", tt.args.val)
		})
	}
}

func TestNewNullTimePtr(t *testing.T) {
	type args struct {
		val time.Time
	}
	tests := []struct {
		name string
		args args
		want *sql.NullTime
	}{
		{
			name: "test",
			args: args{
				val: time.Time{},
			},
			want: &sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewNullTimePtr(tt.args.val), "NewNullTimePtr(%v)", tt.args.val)
		})
	}
}
