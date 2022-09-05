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
