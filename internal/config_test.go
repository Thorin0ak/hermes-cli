package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	c, err := GetConfig()
	assert.Empty(t, err)
	assert.Equal(t, "sse:pxc.dev/123456/", c.Hermes.TopicUri)
}
