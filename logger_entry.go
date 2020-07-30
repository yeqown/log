package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type entry struct {
	logger    *Logger   // logger pointer
	out       io.Writer // write to record
	formatter Formatter // format entry to log
	lv        Level     // the lowest lv which could be log

	fixedField *fixedField            // fixed fields to log
	fields     map[string]interface{} // fields
}

func newEntry(l *Logger) *entry {
	e := entry{
		logger: l,
		out:    l.opt.writer(),
		lv:     l.opt.lv,
		formatter: &TextFormatter{
			isTerminal: l.opt.isTerminal,
		},
		fields: make(Fields, 6),
	}

	if l.opt.globalFields != nil && len(l.opt.globalFields) != 0 {
		copyFields(e.fields, l.opt.globalFields)
	}

	return &e
}

func (e *entry) WithFields(fields Fields) *entry {
	dst := make(Fields, len(fields)+len(e.fields))
	copyFields(dst, fields)
	copyFields(dst, e.fields)

	return &entry{
		logger:    e.logger,
		out:       e.out,
		formatter: e.formatter,
		lv:        e.lv,
		fields:    dst,
	}
}

func (e *entry) reset() {
	e.fields = nil
	e.lv = LevelDebug
	e.out = nil
	e.logger = nil
	e.formatter = nil
	e.fixedField = nil
}

func (e *entry) Fatal(args ...interface{}) {
	e.output(LevelFatal, fmt.Sprint(args...))
	os.Exit(1)
}

func (e *entry) Fatalf(format string, v ...interface{}) {
	e.output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (e *entry) Error(args ...interface{}) {
	e.output(LevelError, fmt.Sprint(args...))
}

func (e *entry) Errorf(format string, v ...interface{}) {
	e.output(LevelError, fmt.Sprintf(format, v...))
}

func (e *entry) Warn(args ...interface{}) {
	e.output(LevelWarning, fmt.Sprint(args...))
}

func (e *entry) Warnf(format string, v ...interface{}) {
	e.output(LevelWarning, fmt.Sprintf(format, v...))
}

func (e *entry) Info(args ...interface{}) {
	e.output(LevelInfo, fmt.Sprint(args...))
}

func (e *entry) Infof(format string, v ...interface{}) {
	e.output(LevelInfo, fmt.Sprintf(format, v...))
}

func (e *entry) Debug(args ...interface{}) {
	e.output(LevelDebug, fmt.Sprint(args...))
}

func (e *entry) Debugf(format string, v ...interface{}) {
	e.output(LevelDebug, fmt.Sprintf(format, v...))
}

func (e *entry) output(lv Level, msg string) {
	if e.lv < lv {
		return
	}

	now := time.Now()
	frm := getCaller()
	e.fixedField = &fixedField{
		File:          frm.File + ":" + strconv.Itoa(frm.Line),
		Fn:            frm.Function,
		Timestamp:     now.Unix(),
		FormattedTime: now.Format(time.RFC3339),
	}

	// setting current lv
	e.lv = lv
	e.fields["msg"] = msg

	// format message
	data, err := e.formatter.Format(e)
	if err != nil {
		panic(err) // FIXME: throw error in a way not panic
	}

	// write into writer
	if _, err = e.out.Write(data); err != nil {
		panic(err)
	}
}

// copyFields copy all fields in src to dst
func copyFields(dst, src Fields) {
	for k := range src {
		dst[k] = src[k]
	}
}