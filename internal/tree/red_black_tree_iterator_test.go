package tree

import (
	"math/rand"
	"testing"

	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

func TestIteratorToVisitFullRBTree(t *testing.T) {
	n := 10000
	arr := generateArray(n)
	rbTree := NewRBTree[int, int](ekit.ComparatorRealNumber[int])
	for _, v := range arr {
		assert.Nil(t, rbTree.Add(v, v))
	}

	arrVisit := make([]int, n)
	id := 0
	for iter := rbTree.Begin(); iter.Valid(); iter.Next() {
		pa, err := iter.Get()
		assert.Nil(t, err)
		arrVisit[id] = pa.Key
		assert.Equal(t, id, pa.Key)
		id++
	}
	assert.Equal(t, n, id)
}

func generateArray(n int) []int {
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = i
	}
	rand.Shuffle(n, func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})
	return res
}
