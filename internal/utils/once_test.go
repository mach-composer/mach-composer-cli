package utils

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

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
			onceMap.Get(k).Do(func() {
				doCounter.Add(1)
			})
			totalCounter.Add(1)
		}(k)
	}

	wg.Wait()

	assert.Equal(t, int64(3), doCounter.Load())
	assert.Equal(t, int64(len(keys)), totalCounter.Load())
}
