package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApp(t *testing.T) {
	require.Equal(t, "foo", "foo")
}
