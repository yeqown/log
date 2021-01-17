package log

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

const (
	_FileKey       = "_filepath"
	_FuncName      = "_func"
	_TimestampKey  = "_ts"
	_FormatTimeKey = "_fmt_time"

	// _interfaceFormat the instruction to format interface value.
	_interfaceFormat string = "%+v"
)

// Formatter to format entry fields and other field
type Formatter interface {
	Format(*entry) ([]byte, error)
}

var _ Formatter = &TextFormatter{}

type TextFormatter struct {
	// isTerminal indicates whether the Logger's out is to a terminal.
	isTerminal bool
}

func newTextFormatter(isTerminal bool) Formatter {
	return &TextFormatter{
		isTerminal: isTerminal,
	}
}

// Format entry into log
func (f *TextFormatter) Format(e *entry) ([]byte, error) {
	b := bytes.NewBuffer(nil)

	// write level and colors
	f.printColoredLevel(b, e)

	// write fixed fields
	f.printFixedFields(b, e.fixedField, e.callerReporter, e.formatTime)

	// write fields
	keys := make([]string, 0, len(e.fields))
	for k := range e.fields {
		keys = append(keys, k)
	}
	// sort by keys
	sort.Strings(keys)
	f.printFields(b, keys, e.fields)

	// write a newline flag
	b.WriteString("\n")

	return b.Bytes(), nil
}

// printColoredLevel colored this output
func (f *TextFormatter) printColoredLevel(b *bytes.Buffer, e *entry) {
	val := e.lv.String()
	// 	val := "[" + e.lv.String() + "]"
	if f.isTerminal {
		val = "\033[" + strconv.Itoa(e.lv.Color()) + "m" + val + "\033[0m"
	}
	b.WriteString(val)
}

// printFixedFields
func (f *TextFormatter) printFixedFields(b *bytes.Buffer, fixed *fixedField, printCaller, formatTime bool) {
	if printCaller {
		appendKeyValue(b, _FileKey, fixed.File)
		appendKeyValue(b, _FuncName, fixed.Fn)
	}

	// TODO(@yeqown): maybe need an option to make these two option coexist
	if formatTime {
		appendKeyValue(b, _FormatTimeKey, fixed.FormattedTime)
	} else {
		appendKeyValue(b, _TimestampKey, fixed.Timestamp)
	}
}

// printFields
func (f *TextFormatter) printFields(b *bytes.Buffer, sortedKeys []string, fields Fields) {
	for _, key := range sortedKeys {
		appendKeyValue(b, key, fields[key])
	}
}

func appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	appendValue(b, value)
}

func appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprintf(_interfaceFormat, value)
	}

	b.WriteString(fmt.Sprintf("%q", stringVal))
}
