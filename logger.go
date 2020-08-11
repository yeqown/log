// Package log .
//
// this log is inspired by `https://github.com/silenceper/log` and `https://github.com/sirupsen/logrus`
// 1. can be set to output to file
// 2. log file can be splitted into files day by day, just like `app.20060102.log`
//
package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	lastSplitTimestamp time.Time = time.Now() // last ts when split the file
)

type Logger struct {
	opt *options

	entryPool sync.Pool // entry pool
}

// NewLogger using os.Stdout and LevelDebug to print log
func NewLogger(opts ...LoggerOption) (*Logger, error) {
	in := make([]LoggerOption, 0, len(opts)+1)
	in = append(in, defaultLoggerOption)
	in = append(in, opts...)
	return newLoggerWithOptions(in...)
}

func newLoggerWithOptions(opts ...LoggerOption) (*Logger, error) {
	dst := new(options)

	for _, opt := range opts {
		if err := opt(dst); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}

	l := Logger{
		opt: dst, // options
		entryPool: sync.Pool{
			New: func() interface{} {
				return &entry{}
			},
		},
	}

	return &l, nil
}

func (l *Logger) Fatal(args ...interface{}) {
	e := l.newEntry()
	e.Fatal(args...)
	l.releaseEntry(e)

	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	e := l.newEntry()
	e.Fatalf(format, args...)
	l.releaseEntry(e)

	os.Exit(1)
}

func (l *Logger) Error(args ...interface{}) {
	e := l.newEntry()
	e.Error(args...)
	l.releaseEntry(e)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	e := l.newEntry()
	e.Errorf(format, args...)
	l.releaseEntry(e)
}

func (l *Logger) Warn(args ...interface{}) {
	e := l.newEntry()
	e.Warn(args...)
	l.releaseEntry(e)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	e := l.newEntry()
	e.Warnf(format, args...)
	l.releaseEntry(e)
}

func (l *Logger) Info(args ...interface{}) {
	e := l.newEntry()
	e.Info(args...)
	l.releaseEntry(e)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	e := l.newEntry()
	e.Infof(format, args...)
	l.releaseEntry(e)
}

func (l *Logger) Debug(args ...interface{}) {
	e := l.newEntry()
	e.Debug(args...)
	l.releaseEntry(e)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	e := l.newEntry()
	e.Debugf(format, args...)
	l.releaseEntry(e)
}

func (l *Logger) newEntry() *entry {
	e, ok := l.entryPool.Get().(*entry)
	if ok {
		e.logger = l
		e.lv = l.opt.level()
		e.out = l.opt.writer()
		e.callerReporter = l.opt.callerReporter
		e.fields = make(Fields, 6)
		copyFields(e.fields, l.opt.globalFields)
		e.formatter = &TextFormatter{isTerminal: l.opt.terminal()}
	}

	return newEntry(l)
}

func (l *Logger) releaseEntry(e *entry) {
	e.reset()
	l.entryPool.Put(e)
}

func (l *Logger) SetLogLevel(level Level) {
	l.opt.lv = level
}

func (l *Logger) SetCallerReporter(b bool) {
	l.opt.callerReporter = b
}

func (l *Logger) WithField(key string, value interface{}) *entry {
	e := l.newEntry()
	defer l.releaseEntry(e)

	return e.WithFields(Fields{key: value})
}

func (l *Logger) WithFields(fields Fields) *entry {
	e := l.newEntry()
	defer l.releaseEntry(e)

	return e.WithFields(fields)
}

// open a file to log
// FIXED: wrong permission of folder and file
func open(file string) (fd *os.File, err error) {
	dir, _ := filepath.Split(file)
	if err = os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}

	if fd, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err != nil {
		return nil, err
	}

	return fd, nil
}

// rename file
func rename(dir, filename string) error {
	return os.Rename(
		assembleFilename(dir, filename, true),                  // old name
		assembleFilename(dir, rotateFilename(filename), false), // new name
	)
}

// assembleFilename, if filename if not end of '.log' then append '.log' to filename
//
// example:
// filename = "app"
// rename to: "app.log"
func assembleFilename(dir, filename string, autoSuffix bool) string {
	if autoSuffix && !strings.HasSuffix(filename, ".log") {
		filename = fmt.Sprintf("%s.log", filename)
	}

	return path.Join(dir, filename)
}

// rotateFilename name the old `{filename}` into `{filename}-{date}`
// rotateFilename(`app.log`) => `app.log-20200730`
func rotateFilename(filename string) string {
	date := lastSplitTimestamp.Format("20060102")
	return fmt.Sprintf("%s-%s", filename, date)
}

// shouldSplitByTime judge by current time and lastLogFileDate
// if now is another day from lastLogFileData, then split the log file
func shouldSplitByTime(now time.Time) bool {
	return now.Day() != lastSplitTimestamp.Day()
}
