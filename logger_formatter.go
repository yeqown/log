package log

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"
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

	// sortField represents whether print fields in order of fields'
	// keys lexicographical order.
	sortField bool

	// formatTime means formatter will format timestamp into formatTimeLayout.
	formatTime       bool
	formatTimeLayout string
}

func newTextFormatter(
	isTerminal, sortField, formatTime bool,
	formatTimeLayout string,
) Formatter {
	return &TextFormatter{
		isTerminal:       isTerminal,
		sortField:        sortField,
		formatTime:       formatTime,
		formatTimeLayout: formatTimeLayout,
	}
}

// Format entry into log
func (f *TextFormatter) Format(e *entry) ([]byte, error) {
	b := bytes.NewBuffer(nil)

	// write level and colors
	f.printColoredLevel(b, e)

	// write fixed fields
	f.printFixedFields(b, e.fixedField, e.callerReporter)

	// write fields
	f.printFields(b, e.fields)

	// write a newline flag
	b.WriteString("\n")

	return b.Bytes(), nil
}

// printColoredLevel colored this output
func (f *TextFormatter) printColoredLevel(b *bytes.Buffer, e *entry) {
	s := e.lv.String()
	// 	s := "[" + e.lv.String() + "]"
	if f.isTerminal {
		s = "\033[" + strconv.Itoa(e.lv.Color()) + "m" + s + "\033[0m"
	}
	b.WriteString(s)
}

// printFixedFields
func (f *TextFormatter) printFixedFields(b *bytes.Buffer, fixed *fixedField, printCaller bool) {
	if printCaller {
		appendKeyValue(b, _FileKey, fixed.File)
		appendKeyValue(b, _FuncName, fixed.Fn)
	}

	// DONE(@yeqown): maybe need an option to make these two option coexist:
	// use WithTimeFormat option API.
	if f.formatTime {
		appendKeyValue(b, _FormatTimeKey,
			time.Unix(fixed.Timestamp, 0).Format(f.formatTimeLayout))
	} else {
		appendKeyValue(b, _TimestampKey, fixed.Timestamp)
	}
}

// printFields append fields into buffer, sortField represents join
// fields in order or not, the order is keys' lexicographical order.
func (f *TextFormatter) printFields(b *bytes.Buffer, fields Fields) {
	if !f.sortField {
		for key := range fields {
			appendKeyValue(b, key, fields[key])
		}
		return
	}

	// If the  formatter need sort keys: WithSortFields option API.
	// sort keys firstly.
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// join append field by order of sorted keys.
	for _, key := range keys {
		appendKeyValue(b, key, fields[key])
	}
	return
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
