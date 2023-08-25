package utils

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var lock = sync.Mutex{}

func TestOnceMapGet(t *testing.T) {
	doCounter := atomic.Int64{}
	totalCounter := atomic.Int64{}

	keys := []string{
		"first",
		"second",
		"first",
		"first",
		"second",
		"third",
	}

	onceMap := OnceMap[string]{}

	wg := sync.WaitGroup{}

	for _, k := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			lock.Lock()
			onceMap.Get(k).Do(func() {
				time.Sleep(time.Second)
				doCounter.Add(1)
			})
			lock.Unlock()
			totalCounter.Add(1)
		}(k)
	}

	wg.Wait()

	assert.Equal(t, int64(3), doCounter.Load())
	assert.Equal(t, int64(len(keys)), totalCounter.Load())
}
