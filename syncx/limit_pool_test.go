package syncx

import (
	"bytes"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitPool(t *testing.T) {

	expectedMaxAttempts := 3
	expectedVal := []byte("A")

	pool := NewLimitPool(expectedMaxAttempts, func() []byte {
		var buffer bytes.Buffer
		buffer.Write(expectedVal)
		return buffer.Bytes()
	})

	var wg sync.WaitGroup
	bufChan := make(chan []byte, expectedMaxAttempts)
	for i := 0; i < expectedMaxAttempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := pool.Get()
			assert.NotZero(t, buf)
			assert.Equal(t, string(expectedVal), string(buf))
			bufChan <- buf
		}()
	}

	wg.Wait()
	close(bufChan)

	// 超过最大申请次数返回零值
	assert.Zero(t, pool.Get())

	// 归还一个
	pool.Put(<-bufChan)

	// 再次申请仍可以拿到非零值缓冲区
	assert.NotZero(t, string(expectedVal), string(pool.Get()))
}
