package cli

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun_Default(t *testing.T) {
	h := NewTerraformHook("test")
	msg := "this is a test message"

	e := &zerolog.Event{}
	e.Ctx(ContextWithOutput(context.Background(), OutputTypeConsole))

	sink := &LogSink{}

	defer SetSinkLogger(sink)()

	h.Run(e, zerolog.InfoLevel, msg)

	assert.Equal(t, "info", sink.Index(0).Level)
	assert.Equal(t, msg, sink.Index(0).Message)
}

func TestRun_WithJson(t *testing.T) {
	h := NewTerraformHook("test")
	msg := `{"@level":"info","@message":"Terraform 1.7.3","@module":"terraform.ui","@timestamp":"2024-02-23T15:50:09.437498+01:00","terraform":"1.7.3","type":"version","ui":"1.2"}
{"@level":"info","@message":"Plan: 0 to add, 0 to change, 0 to destroy.","@module":"terraform.ui","@timestamp":"2024-02-23T15:50:09.439546+01:00","changes":{"add":0,"change":0,"import":0,"remove":0,"operation":"plan"},"type":"change_summary"}
{"@level":"info","@message":"Outputs: 1","@module":"terraform.ui","@timestamp":"2024-02-23T15:50:09.439577+01:00","outputs":{"hash":{"sensitive":false,"action":"create"}},"type":"outputs"}
{"@level":"info","@message":"Apply complete! Resources: 0 added, 0 changed, 0 destroyed.","@module":"terraform.ui","@timestamp":"2024-02-23T15:50:09.444929+01:00","changes":{"add":0,"change":0,"import":0,"remove":0,"operation":"apply"},"type":"change_summary"}
{"@level":"warn","@message":"Outputs: 1","@module":"terraform.ui","@timestamp":"2024-02-23T15:50:09.444947+01:00","outputs":{"hash":{"sensitive":false,"type":"string","value":"74234e98afe7498fb5daf1f36ac2d78acc339464f950703b8c019892f982b90b"}},"type":"outputs"}`

	e := &zerolog.Event{}
	e.Ctx(ContextWithOutput(context.Background(), OutputTypeJSON))

	sink := &LogSink{}

	defer SetSinkLogger(sink)()

	h.Run(e, zerolog.InfoLevel, msg)

	assert.Equal(t, 5, len(sink.logs))
	assert.Equal(t, "info", sink.Index(0).Level)
	assert.Equal(t, "Terraform 1.7.3", sink.Index(0).Message)
	assert.Equal(t, "warn", sink.Index(4).Level)
}
