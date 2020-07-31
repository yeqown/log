package log

import (
	"strconv"
	"sync"
	"testing"
	"time"

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
	filename := "split.log"

	// prepare file
	_, err := open(assembleFilename(dir, filename, true))
	assert.Nil(t, err)

	// rename file
	err = rename(dir, filename)
	assert.Nil(t, err)

	// renew file
	_, err = open(assembleFilename(dir, filename, true))
	assert.Nil(t, err)
}

func Test_Logger_FileSplit(t *testing.T) {
	dir := "./testdata"
	filename := strconv.FormatInt(time.Now().Unix(), 10) + ".log"

	l, err := NewLogger(
		WithFileLog(assembleFilename(dir, filename, true), true),
		WithStdout(true),
	)
	assert.Nil(t, err)

	// ticker := time.NewTicker(1 * time.Second)
	threshold := 3
	once := sync.Once{}

	switchWriter := func() {
		t.Log("switch writer")

		// rename file
		err = rename(dir, filename)
		assert.Nil(t, err)

		// renew file
		fd, err := open(assembleFilename(dir, filename, true))
		assert.Nil(t, err)
		l.opt.w = fd
	}

	for counter := 1; counter < 100; counter++ {
		if counter > 6 {
			break
		}
		if counter > threshold {
			once.Do(func() {
				switchWriter()
			})
		}

		l.Infof("count=%d", counter)
	}
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

func Test_shouldSplitByTime(t *testing.T) {
	now := time.Now()

	type args struct {
		lastSplitTimestamp time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 0",
			args: args{
				lastSplitTimestamp: now, // now
			},
			want: false,
		},
		{
			name: "case 1",
			args: args{
				lastSplitTimestamp: now.Add(-24 * time.Hour), // yesterday
			},
			want: true,
		},
		{
			name: "case 2",
			args: args{
				lastSplitTimestamp: now.Add(-48 * time.Hour),
			}, // the day before yesterday
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set var
			lastSplitTimestamp = tt.args.lastSplitTimestamp

			if got := shouldSplitByTime(now); got != tt.want {
				t.Errorf("shouldSplitByTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
