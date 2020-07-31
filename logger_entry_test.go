package log

import (
	"bytes"
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
	assert.Contains(t, entry2.fields, "foo")
	assert.Equal(t, "bar updated", entry2.fields["foo"])
}
