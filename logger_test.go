package log

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLogger(t *testing.T) {
	logger, err := NewLogger()
	assert.Nil(t, err)

	logger.Error("Error")
	logger.Warn("Warn")
	logger.Info("Info")
	logger.Debug("Debug")
}

func Test_FileSplit(t *testing.T) {
	dir := "./testdata"
	filename := "split"

	_, err := NewLogger()
	assert.Nil(t, err)

	// rename file
	err = rename(dir, filename)
	assert.Nil(t, err)

	// renew file
	_, err = open(assembleFilename(dir, filename))
	assert.Nil(t, err)
}

func Test_Logger_concurrent(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(10)
	l, _ := NewLogger()

	for i := 0; i < 10; i++ {
		go func(idx int) {
			defer wg.Done()
			l.WithField("key_"+strconv.Itoa(idx), idx).Infof("idx=%d", idx)
		}(i)
	}

	wg.Wait()
}
