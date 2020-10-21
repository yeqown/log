package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newEntry(t *testing.T) {
	b := &bytes.Buffer{}
	l, err := NewLogger(
		WithCustomWriter(b),
		WithGlobalFields(Fields{"foo": "bar"}),
	)
	assert.Nil(t, err)

	entry := newEntry(l)
	assert.Equal(t, l, entry.logger)
	assert.Equal(t, l.opt.level(), entry.lv)
	assert.Equal(t, l.opt.w, entry.out)
	assert.Equal(t, l.opt.globalFields, entry.fields)
	assert.Nil(t, entry.fixedField)
}

func Test_newEntry_Cmp(t *testing.T) {
	b := &bytes.Buffer{}
	l, err := NewLogger(
		WithCustomWriter(b),
		WithGlobalFields(Fields{"foo": "bar"}),
	)
	assert.Nil(t, err)

	// create an entry and release one
	l.WithFields(Fields{"malloc": "release"}).Info("haha")

	entry := newEntry(l)
	entry2 := l.newEntry()

	assert.Equal(t, entry.ctxParser.(funcContextParser).fieldName,
		entry2.ctxParser.(funcContextParser).fieldName) // function couldn't be compared.
	assert.Equal(t, entry.ctx, entry2.ctx)
	assert.Equal(t, entry.lv, entry2.lv)
	assert.Equal(t, entry.fields, entry2.fields)
	assert.Equal(t, entry.callerReporter, entry2.callerReporter)
	assert.Equal(t, entry.formatter, entry2.formatter)
	assert.Equal(t, entry.logger, entry2.logger)
}

func Test_entry_WithFields(t *testing.T) {
	b := &bytes.Buffer{}
	l, err := NewLogger(
		WithCustomWriter(b),
		WithGlobalFields(Fields{"foo": "bar"}),
	)
	assert.Nil(t, err)

	entry := newEntry(l)
	assert.Equal(t, l.opt.globalFields, entry.fields)
	entry2 := entry.WithFields(Fields{
		"foo2": "bar2",
		"foo":  "bar updated",
	})

	assert.Contains(t, entry2.fields, "foo2")
	assert.NotContains(t, entry2.fields, _defaultFieldName)
	assert.Contains(t, entry2.fields, "foo")
	assert.Equal(t, "bar updated", entry2.fields["foo"])
}

func Test_entry_Without_Caller(t *testing.T) {
	b := &bytes.Buffer{}
	l, err := NewLogger(
		WithCustomWriter(b),
	)
	assert.Nil(t, err)

	entry := newEntry(l)
	assert.Equal(t, false, entry.callerReporter)
	entry.Info("with out caller")
	assert.NotContains(t, b.String(), "file")
	assert.NotContains(t, b.String(), "fn")

	b.Reset()

	// open
	l.SetCallerReporter(true)
	entry2 := newEntry(l)
	assert.Equal(t, true, entry2.callerReporter)
	entry2.Info("with caller")
	assert.Contains(t, b.String(), "file")
	assert.Contains(t, b.String(), "fn")
}

func Test_entry_WithContextAndWithFields(t *testing.T) {
	b := &bytes.Buffer{}
	l, err := NewLogger(
		WithCustomWriter(b),
	)
	assert.Nil(t, err)

	//  withContext then withFields
	entry2 := l.newEntry().WithContext(context.TODO()).WithFields(Fields{
		"field1": "field1",
	})
	entry2.Info("output")

	assert.Contains(t, entry2.fields, _defaultFieldName)
	assert.Contains(t, entry2.fields, "field1")
	assert.Equal(t, "non action", entry2.fields[_defaultFieldName])

	// withFields then withContext
	entry3 := l.newEntry().WithFields(Fields{
		"field2": "field2",
	}).WithContext(context.TODO())
	entry3.Info("output")

	assert.Contains(t, entry3.fields, _defaultFieldName)
	assert.Contains(t, entry3.fields, "field2")

}

func Test_entry_WithContextButNotSetParser(t *testing.T) {
	// b := &bytes.Buffer{}
	l, err := NewLogger(
	// WithCustomWriter(b),
	)
	assert.Nil(t, err)

	entry := l.newEntry()
	entry2 := entry.WithContext(context.TODO())
	entry2.Info("output")

	assert.Contains(t, entry2.fields, _defaultFieldName)
	assert.Equal(t, "non action", entry2.fields[_defaultFieldName])
}

func Test_entry_WithContextAndSetParser(t *testing.T) {
	var customParser = func(ctx context.Context) interface{} {
		return "custom"
	}

	// b := &bytes.Buffer{}
	l, err := NewLogger(
		// WithCustomWriter(b),
		WithContextParser(NewContextParserFunc(customParser, "ctxField")),
	)
	assert.Nil(t, err)

	entry := l.newEntry()
	entry2 := entry.WithContext(context.TODO())
	entry2.Info("output")

	assert.Contains(t, entry2.fields, "ctxField")
	assert.Equal(t, "custom", entry2.fields["ctxField"])
}
