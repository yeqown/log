// Package log
//
// this log is based `https://github.com/silenceper/log` but more functions:
//
// 1. can be set to output to file
//
// 2. log file can splited into files day by day, just like `app.20060102.log`
//
package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type (
	Level int
)

const (
	LevelFatal   = iota // FatalLevel
	LevelError          // ErrorLevel
	LevelWarning        // WarningLevel
	LevelInfo           // InfoLevel
	LevelDebug          // DebugLevel
)

var (
	_log           *logger   = NewLogger() // default logger
	lstLogFileDate time.Time = time.Now()  // last date time when split logfile
)

func Fatal(args ...interface{}) {
	_log.Output(LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	_log.Output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Error(args ...interface{}) {
	_log.Output(LevelError, fmt.Sprint(args...))
}

func Errorf(format string, v ...interface{}) {
	_log.Output(LevelError, fmt.Sprintf(format, v...))
}

func Warn(args ...interface{}) {
	_log.Output(LevelWarning, fmt.Sprint(args...))
}

func Warnf(format string, v ...interface{}) {
	_log.Output(LevelWarning, fmt.Sprintf(format, v...))
}

func Info(args ...interface{}) {
	_log.Output(LevelInfo, fmt.Sprint(args...))
}

func Infof(format string, v ...interface{}) {
	_log.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func Debug(args ...interface{}) {
	_log.Output(LevelDebug, fmt.Sprint(args...))
}

func Debugf(format string, v ...interface{}) {
	_log.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func SetLogLevel(level Level) {
	_log.SetLogLevel(level)
}

func SetFileOutput(logPath, filename string) {
	_log.SetFileOutput(logPath, filename)
}

type logger struct {
	stdLog   *log.Logger // os.stderr
	fileLog  *log.Logger // 文件
	logLevel Level       // 小于等于该级别的level才会被记录
}

// NewLogger 实例化，供自定义
func NewLogger() *logger {
	return &logger{
		stdLog:   log.New(os.Stderr, "", log.Llongfile|log.LstdFlags),
		fileLog:  nil,
		logLevel: LevelDebug,
	}
}

// SetFileOutput to set file output and create new croutine
// to recv time.Ticker with 1 min interval.
func (l *logger) SetFileOutput(logPath, filename string) {
	file := openOrCreate(assembleFilepath(logPath, filename))
	l.fileLog = log.New(file, filename, log.Llongfile|log.LstdFlags)

	// new croutine to split file
	go func(logPath, filename string) {
		ticker := time.NewTicker(1 * time.Minute)
		for true {
			select {
			case <-ticker.C:
				if !timeToSplit() {
					continue
				}
				renameLogfile(logPath, filename)
				file := openOrCreate(assembleFilepath(logPath, filename))
				l.fileLog = log.New(file, "", log.Lshortfile|log.LstdFlags)
				lstLogFileDate = time.Now()
			}
		}
	}(logPath, filename)
}

// the most based function to log
func (l *logger) Output(level Level, s string) {
	if l.logLevel < level {
		return
	}
	formatStr := "[UNKNOWN] %s"
	formatFileStr := "[UNKNOWN] %s"
	switch level {
	case LevelFatal:
		formatStr = "\033[35m[FATAL]\033[0m %s"
		formatFileStr = "[FATAL] %s"
	case LevelError:
		formatStr = "\033[31m[ERROR]\033[0m %s"
		formatFileStr = "[ERROR] %s"
	case LevelWarning:
		formatStr = "\033[33m[WARN]\033[0m %s"
		formatFileStr = "[WARN] %s"
	case LevelInfo:
		formatStr = "\033[32m[INFO]\033[0m %s"
		formatFileStr = "[INFO] %s"
	case LevelDebug:
		formatStr = "\033[36m[DEBUG]\033[0m %s"
		formatFileStr = "[DEBUG] %s"
	}
	stdFormat := fmt.Sprintf(formatStr, s)
	fileFormat := fmt.Sprintf(formatFileStr, s)

	file, function, line := findCaller(5)
	println(file, function, line)

	// output to os.stderr
	if err := l.stdLog.Output(3, stdFormat); err != nil {
		panic(err)
	}

	// output to file
	if l.fileLog == nil {
		return
	}

	if err := l.fileLog.Output(3, fileFormat); err != nil {
		panic(err)
	}
}

func (l *logger) Fatal(args ...interface{}) {
	l.Output(LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *logger) Error(args ...interface{}) {
	l.Output(LevelError, fmt.Sprint(args...))
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *logger) Warn(args ...interface{}) {
	l.Output(LevelWarning, fmt.Sprint(args...))
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Output(LevelWarning, fmt.Sprintf(format, v...))
}

func (l *logger) Info(args ...interface{}) {
	l.Output(LevelInfo, fmt.Sprint(args...))
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *logger) Debug(args ...interface{}) {
	l.Output(LevelDebug, fmt.Sprint(args...))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *logger) SetLogLevel(level Level) {
	l.logLevel = level
}

//////////////////////////////////////
//
// logger utils functions
//
//////////////////////////////////////
func openOrCreate(filepath string) *os.File {
	file := new(os.File)
	var err error
	if file, err = os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644); err != nil {
		panic(err)
	}
	return file
}

func assembleFilepath(logPath, filename string) string {
	// if not end of `.log` append filename
	// example:
	//
	// filename = "app"
	// rename to: "app.log"
	//
	if !strings.HasSuffix(filename, ".log") {
		filename = fmt.Sprintf("%s.log", filename)
	}
	return path.Join(logPath, filename)
}

func formatFilename(filename string) string {
	date := lstLogFileDate.Format("20060102")
	return fmt.Sprintf("%s-%s", filename, date)
}

func renameLogfile(logPath, filename string) {
	if err := os.Rename(
		assembleFilepath(logPath, filename),
		assembleFilepath(logPath, formatFilename(filename)),
	); err != nil {
		panic(err)
	}
}

func timeToSplit() bool {
	now := time.Now()
	if now.Day() != lstLogFileDate.Day() {
		return true
	}
	return false
}
