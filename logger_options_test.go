package log

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isTerminal(t *testing.T) {
	var w io.Writer = os.Stdout
	ok := isTerminal(w)
	assert.True(t, ok)

	buf := bytes.NewBuffer(nil)
	w = buf
	ok = isTerminal(w)
	assert.False(t, ok)

	fd, err := os.OpenFile("./testdata/is_bytes_outputing_device.t", os.O_CREATE|os.O_TRUNC, 0666)
	assert.NoError(t, err)
	w = fd
	ok = isTerminal(w)
	assert.False(t, ok)

	w = io.MultiWriter(os.Stdout, fd)
	ok = isTerminal(w)
	assert.False(t, ok)
}
