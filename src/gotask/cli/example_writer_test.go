package cli

import (
	"bytes"
	"strings"

	"github.com/bmizerany/assert"

	"testing"
)

func TestExampleWriter_Write(t *testing.T) {
	var out bytes.Buffer
	b := exampleWriter{
		"foo",
	}
	b.Write(&out)

	assert.Tf(t, strings.Contains(out.String(), `package foo`), "%v", out.String())
}
