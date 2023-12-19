package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"testing"
)

func TestParseHashOutputNoHashAttribute(t *testing.T) {
	val := cty.ObjectVal(map[string]cty.Value{})

	_, err := ParseHashOutput(val)
	assert.ErrorContains(t, err, "no attribute with key hash found in terraform output")
}

func TestParseHashOutputInvalidResponse(t *testing.T) {
	val := cty.ObjectVal(map[string]cty.Value{
		"hash": cty.StringVal("invalid"),
	})

	_, err := ParseHashOutput(val)
	assert.ErrorAs(t, err, &cty.PathError{})
}

func TestParseHashOutputNilResponse(t *testing.T) {
	val := cty.ObjectVal(map[string]cty.Value{
		"hash": cty.ObjectVal(map[string]cty.Value{
			"sensitive": cty.BoolVal(false),
			"value":     cty.NilVal,
			"type":      cty.StringVal("some-type"),
		}),
	})

	_, err := ParseHashOutput(val)
	assert.ErrorContains(t, err, "no value set for hash in terraform output")
}

func TestParseHashOutputOk(t *testing.T) {
	hash := "1234567890"
	val := cty.ObjectVal(map[string]cty.Value{
		"hash": cty.ObjectVal(map[string]cty.Value{
			"sensitive": cty.BoolVal(false),
			"value":     cty.StringVal(hash),
			"type":      cty.StringVal("some-type"),
		}),
	})

	parsedHash, err := ParseHashOutput(val)
	assert.NoError(t, err)
	assert.Equal(t, hash, parsedHash)
}
