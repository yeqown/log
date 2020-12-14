package log

import (
	"bytes"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewLogger_WithWriter(t *testing.T) {
	b := &bytes.Buffer{}
	logger, err := NewLogger(
		WithCustomWriter(b),
	)
	assert.Nil(t, err)

	logger.
		WithField("key1", "value1").
		WithFields(Fields{"key2": "value2"}).
		Info("fields test")
	assert.Contains(t, b.String(), "key1")
	assert.Contains(t, b.String(), "value1")
	assert.Contains(t, b.String(), "key2")
	assert.Contains(t, b.String(), "value2")
}

func Test_NewLogger_WithGlobalFields(t *testing.T) {
	b := &bytes.Buffer{}
	logger, err := NewLogger(
		WithCustomWriter(b),
		WithGlobalFields(Fields{"global_key": "global_value"}),
	)
	assert.Nil(t, err)

	// case 0
	logger.
		WithField("key1", "value2222").Info("fields test")
	assert.Contains(t, b.String(), "value2222")
	assert.Contains(t, b.String(), "global_value")

	b.Reset()
	assert.NotContains(t, b.String(), "value2222")
	assert.NotContains(t, b.String(), "global_value")

	// case 1 if duplicated the key with global key, then replace the global key
	logger.WithFields(Fields{
		"global_key": "global_valueless",
		"key2":       "value2222",
	}).Info("fields test")
	t.Log(b.String())
	assert.Contains(t, b.String(), "global_valueless")
	assert.Contains(t, b.String(), "value2222")
}

func Test_Logger_SetLevel(t *testing.T) {
	b := &bytes.Buffer{}
	logger, err := NewLogger(
		WithCustomWriter(b),
	)
	assert.Nil(t, err)

	// set level
	logger.SetLogLevel(LevelWarning)

	logger.Debug("debug")
	assert.NotContains(t, b.String(), "DBG")
	logger.Info("info")
	assert.NotContains(t, b.String(), "INF")
	logger.Warn("warn")
	assert.Contains(t, b.String(), "WRN")
	logger.Error("error")
	assert.Contains(t, b.String(), "ERR")
}

func Test_Logger_WithLevel(t *testing.T) {
	b := &bytes.Buffer{}
	logger, err := NewLogger(
		WithCustomWriter(b),
		WithLevel(LevelWarning),
	)
	assert.Nil(t, err)

	logger.Debug("debug")
	assert.NotContains(t, b.String(), "DBG")
	logger.Info("info")
	assert.NotContains(t, b.String(), "INF")
	logger.Warn("warn")
	assert.Contains(t, b.String(), "WRN")
	logger.Error("error")
	assert.Contains(t, b.String(), "ERR")
}

func Test_Logger_SetTimeFormat(t *testing.T) {
	b := &bytes.Buffer{}
	logger, err := NewLogger(
		WithCustomWriter(b),
	)
	assert.Nil(t, err)
	logger.SetTimeFormat(true, "")

	logger.Info("info test time format")
	assert.Contains(t, b.String(), _FormatTimeKey)
	assert.NotContains(t, b.String(), _TimestampKey)
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
