package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	c := GetConfig()
	assert.NotEmpty(t, c)
	assert.Equal(t, "sse:pxc.dev/123456/", c.Hermes.TopicUri)
}
