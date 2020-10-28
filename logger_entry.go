package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type entry struct {
	logger         *Logger   // logger pointer
	out            io.Writer // write to record
	formatter      Formatter // format entry to log
	lv             Level     // the lowest lv which could be log
	callerReporter bool      // log caller

	fixedField *fixedField // fixed fields to log
	fields     Fields      // fields

	ctx       context.Context
	ctxParser ContextParser
}

func newEntry(l *Logger) *entry {
	e := entry{
		logger: l,
		out:    l.opt.writer(),
		lv:     l.opt.lv,
		formatter: &TextFormatter{
			isTerminal: l.opt.isTerminal,
		},
		callerReporter: l.opt.callerReporter,
		fixedField:     nil,
		fields:         make(Fields, 4),
		ctxParser:      l.opt.ctxParser,
	}

	if l.opt.globalFields != nil && len(l.opt.globalFields) != 0 {
		copyFields(e.fields, l.opt.globalFields)
	}

	return &e
}

func (e *entry) WithFields(fields Fields) *entry {
	dst := make(Fields, len(fields)+len(e.fields))
	// FIXED: copy entry's fields at first, then copy newer fields
	copyFields(dst, e.fields)
	copyFields(dst, fields)

	return &entry{
		logger:         e.logger,
		out:            e.out,
		formatter:      e.formatter,
		lv:             e.lv,
		fields:         dst,
		ctx:            e.ctx,
		ctxParser:      e.ctxParser,
		callerReporter: e.callerReporter,
	}
}

// WithContext would overwrite the previous ctx which exists in `e`.
func (e *entry) WithContext(ctx context.Context) *entry {
	return &entry{
		logger:         e.logger,
		out:            e.out,
		formatter:      e.formatter,
		lv:             e.lv,
		fields:         e.fields, // TODO: ANY PROBLEM HERE ?
		ctx:            ctx,
		ctxParser:      e.ctxParser,
		callerReporter: e.callerReporter,
	}
}

func (e *entry) reset() {
	e.fields = nil
	e.lv = LevelDebug
	e.out = nil
	e.logger = nil
	e.formatter = nil
	e.fixedField = nil
	e.ctx = nil
	e.ctxParser = nil
	e.callerReporter = false
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

	e.fixedField = &fixedField{
		//File:          file + ":" + strconv.Itoa(line),
		//Fn:            fn,
		Timestamp:     now.Unix(),
		FormattedTime: now.Format(time.RFC3339),
	}

	if e.callerReporter {
		file := "failed"
		fn := "failed"
		line := 0

		frm := getCaller()
		if frm != nil {
			file = frm.File
			fn = frm.Function
			line = frm.Line
		}

		e.fixedField.File = file + ":" + strconv.Itoa(line)
		e.fixedField.Fn = fn
	}

	// setting current lv
	e.lv = lv
	e.fields["msg"] = msg

	// parse context
	if e.ctx != nil && e.ctxParser != nil {
		_ctxValue := e.ctxParser.Parse(e.ctx)
		e.fields[e.ctxParser.FieldName()] = _ctxValue
	}

	// format message
	data, err := e.formatter.Format(e)
	if err != nil {
		// FIXED: throw error in a way not panic
		// panic(err)
		log.Printf("WARN: could not format message, err=%v", err)
	}

	// write into writer
	if _, err = e.out.Write(data); err != nil {
		// panic(err)
		log.Printf("WARN: could not write log data, err=%v", err)
	}
}
