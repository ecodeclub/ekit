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

package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeEqual(t *testing.T) {
	tests := []struct {
		name  string
		equal EqualFunc[any]

		want bool
	}{
		{
			name: "panic",
			equal: func(x, y any) bool {
				panic("panic test")
			},

			want: true,
		},
		{
			name: "no panic",
			equal: func(x, y any) bool {
				return true
			},

			want: false,
		},
	}
	for _, tt := range tests {
		isPanic, _ := tt.equal.safeEqual(1, 1)
		assert.Equal(t, tt.want, isPanic)
	}

}
