/*
 * this log is based `https://github.com/silenceper/log`
 * just more functions:
 *
 * 1. can be set to output to file
 * 2. log file can splited into files day by day, just like `app.20060102.log`
 *
 * if want to log to file you must NewLogger and SetFileOutput
 */
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
	LevelFatal = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

var (
	_log           *logger   = New()
	lstLogFileDate time.Time = time.Now()
)

func Fatal(s string) {
	_log.Output(LevelFatal, s)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	_log.Output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Error(s string) {
	_log.Output(LevelError, s)
}

func Errorf(format string, v ...interface{}) {
	_log.Output(LevelError, fmt.Sprintf(format, v...))
}

func Warn(s string) {
	_log.Output(LevelWarning, s)
}

func Warnf(format string, v ...interface{}) {
	_log.Output(LevelWarning, fmt.Sprintf(format, v...))
}

func Info(s string) {
	_log.Output(LevelInfo, s)
}

func Infof(format string, v ...interface{}) {
	_log.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func Debug(s string) {
	_log.Output(LevelDebug, s)
}

func Debugf(format string, v ...interface{}) {
	_log.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func SetLogLevel(level Level) {
	_log.SetLogLevel(level)
}

type logger struct {
	std_log  *log.Logger
	file_log *log.Logger
	//小于等于该级别的level才会被记录
	logLevel Level
}

//NewLogger 实例化，供自定义
func NewLogger() *logger {
	return &logger{
		std_log:  log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags),
		file_log: nil,
		logLevel: LevelDebug,
	}
}

//New 实例化，供外部直接调用 log.XXXX
func New() *logger {
	return &logger{
		std_log:  log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags),
		file_log: nil,
		logLevel: LevelDebug,
	}
}

/*
 * @logPath string
 * @filename string
 */
func (l *logger) SetFileOutput(logPath, filename string) error {
	file, err := openOrCreate(assembleFilepath(logPath, filename))
	if err != nil {
		return err
	}
	l.file_log = log.New(file, "", log.Lshortfile|log.LstdFlags)

	// new croutine to split file
	go func(logPath, filename string) {
		for true {
			now := time.Now()
			if shouldSplit() {
				// rename old file
				if err := renameLogfile(logPath, filename); err != nil {
					panic(err)
				}
				// renew file
				if file, err := openOrCreate(assembleFilepath(logPath, filename)); err != nil {
					panic(err)
				} else {
					l.file_log = log.New(file, "", log.Lshortfile|log.LstdFlags)
					lstLogFileDate = now
				}
			}
			time.Sleep(60 * time.Second)
		}
	}(logPath, filename)
	return nil
}

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
	std_s := fmt.Sprintf(formatStr, s)
	file_s := fmt.Sprintf(formatFileStr, s)
	// output to stderr
	if err := l.std_log.Output(3, std_s); err != nil {
		panic(err)
	}
	// output to file
	if l.file_log != nil {
		if err := l.file_log.Output(3, file_s); err != nil {
			panic(err)
		}
	}
}

func (l *logger) Fatal(s string) {
	l.Output(LevelFatal, s)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *logger) Error(s string) {
	l.Output(LevelError, s)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *logger) Warn(s string) {
	l.Output(LevelWarning, s)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Output(LevelWarning, fmt.Sprintf(format, v...))
}

func (l *logger) Info(s string) {
	l.Output(LevelInfo, s)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *logger) Debug(s string) {
	l.Output(LevelDebug, s)
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

func openOrCreate(filepath string) (file *os.File, err error) {
	if file, err = os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644); err == nil {
		return file, nil
	}
	return nil, err
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

func renameLogfile(logPath, filename string) error {
	if err := os.Rename(
		assembleFilepath(logPath, filename),
		assembleFilepath(logPath, formatFilename(filename)),
	); err != nil {
		return err
	}
	return nil
}

func shouldSplit() bool {
	now := time.Now()
	if now.Day() != lstLogFileDate.Day() {
		return true
	}
	return false
}
