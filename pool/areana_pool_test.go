package pool

import (
	"arena"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type Test struct {
	A int
}

func TestSyncPool(t *testing.T) {
	pool := sync.Pool{New: func() any {
		return &Test{
			A: 1,
		}
	},
	}
	testObject := pool.Get().(*Test)
	assert.Equal(t, testObject.A, 1)
	pool.Put(testObject)

}

func TestArena(t *testing.T) {
	mem := arena.NewArena()

	obj := arena.New[Test](mem)
	mem.Free()
	obj.A = 2

	fmt.Println(obj)

	//slice := arena.MakeSlice[Test](mem, 100, 200)
	//fmt.Println(slice)
}

func TestArenaPool(t *testing.T) {
	pool := NewArenaPool(func() any {
		return &Test{
			A: 1,
		}
	})
	testObject, err := pool.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, testObject.Object.(*Test).A)
	pool.Put(testObject)
	testObject1, err := pool.Get()
	assert.Equal(t, testObject, testObject1)

	testObject1.Free()

}
