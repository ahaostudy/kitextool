package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	require.NoError(t, err)
	require.Equal(t, "192.168", ip[:7])
}
