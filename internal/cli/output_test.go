package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertOutputType_Invalid(t *testing.T) {
	_, err := ConvertOutputType("invalid")
	assert.Error(t, err)
}

func TestConvertOutputType_OK(t *testing.T) {
	o, err := ConvertOutputType("console")
	assert.NoError(t, err)
	assert.Equal(t, OutputTypeConsole, o)
}
