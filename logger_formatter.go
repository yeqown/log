package log

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

// Formatter to format entry fields and other field
type Formatter interface {
	Format(*entry) ([]byte, error)
}

var _ Formatter = &TextFormatter{}

type TextFormatter struct {
	// Whether the Logger's out is to a terminal
	isTerminal bool
}

// Format entry into log
func (f *TextFormatter) Format(e *entry) ([]byte, error) {
	b := bytes.NewBuffer(nil)

	// write level and colors
	f.printColoredLevel(b, e)

	// write fixed fields
	f.printFixedFields(b, e.fixedField, e.callerReporter)

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

// colored this output
func (f *TextFormatter) printColoredLevel(b *bytes.Buffer, e *entry) {
	val := "[" + e.lv.String() + "]"
	if f.isTerminal {
		val = "\033[" + strconv.Itoa(e.lv.Color()) + "m[" + e.lv.String() + "]\033[0m"
	}
	b.WriteString(val)
}

// printFixedFields
func (f *TextFormatter) printFixedFields(b *bytes.Buffer, fixed *fixedField, printCaller bool) {
	if printCaller {
		appendKeyValue(b, "file", fixed.File)
		appendKeyValue(b, "fn", fixed.Fn)
	}
	appendKeyValue(b, "timestamp", fixed.Timestamp)
	appendKeyValue(b, "formatted_time", fixed.FormattedTime)
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
		stringVal = fmt.Sprint(value)
	}

	b.WriteString(fmt.Sprintf("%q", stringVal))
}
