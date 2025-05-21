package log

import (
	"bytes"
	"fmt"
	"sort"
	"time"
)

const (
	_FileKey     = "_file"
	_FuncNameKey = "_func"
	// _TimestampKey  = "_timestamp"
	// _FormatTimeKey = "_time"

	// _interfaceFormat the instruction to format interface value.
	_interfaceFormat string = "%+v"
)

// Formatter to format entry fields and other field
type Formatter interface {
	Format(*entry, string) ([]byte, error)
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
func (f *TextFormatter) Format(e *entry, msg string) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	// write level and colors
	f.printColoredLevel(b, e)
	// write fixed fields
	f.printFixedFields(b, e.fixedField, e.withCaller)
	// write fields
	if len(e.fields) > 0 {
		f.printFields(b, e.fields)
	}
	// write a newline flag
	b.WriteString(" " + msg + "\n")

	return b.Bytes(), nil
}

// printColoredLevel colored this output
func (f *TextFormatter) printColoredLevel(b *bytes.Buffer, e *entry) {
	// s := e.lv.String()
	s := "[" + e.lv.String() + "]"
	if f.isTerminal {
		s = "\033[" + e.lv.Color() + "m" + s + "\033[0m"
	}
	s += " "
	b.WriteString(s)
}

// printFixedFields
func (f *TextFormatter) printFixedFields(b *bytes.Buffer, fixed *fixedField, printCaller bool) {
	// DONE(@yeqown): maybe need an option to make these two option coexist:
	// use WithTimeFormat option API.
	if f.formatTime {
		appendValue(b, time.Unix(fixed.Timestamp, 0).Format(f.formatTimeLayout), false)
	} else {
		appendValue(b, fixed.Timestamp, false)
	}

	if printCaller {
		appendKeyValue(b, _FileKey, fixed.File, true, false)
		appendKeyValue(b, _FuncNameKey, fixed.Fn, true, false)
	}
}

// printFields append fields into buffer, sortField represents join
// fields in order or not, the order is keys' lexicographical order.
func (f *TextFormatter) printFields(b *bytes.Buffer, fields Fields) {
	b.WriteString(" Fields{")
	defer b.WriteString("}")

	if !f.sortField {
		n := 0
		for key := range fields {
			appendKeyValue(b, key, fields[key], n != 0, true)
			n++
		}
		return
	}

	// If the formatter needs sort keys: WithSortFields option API.
	// Sort keys firstly.
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// join the appended field by order of sorted keys.
	n := 0
	for _, key := range keys {
		appendKeyValue(b, key, fields[key], n != 0, true)
		n++
	}
}

func appendKeyValue(b *bytes.Buffer, key string, value interface{}, indent, withQuote bool) {
	if b.Len() > 0 && indent {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	appendValue(b, value, withQuote)
}

func appendValue(b *bytes.Buffer, value interface{}, withQuote bool) {
	stringVal, ok := value.(string)
	if ok {
		if !withQuote {
			b.WriteString(stringVal)
			return
		}
	}

	stringVal = fmt.Sprintf(_interfaceFormat, value)
	b.WriteString(fmt.Sprintf("%q", stringVal))
}
