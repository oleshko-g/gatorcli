package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfigFilePath(t *testing.T) {
	configFilePath, err := getConfigFilePath()
	require.NoError(t, err)
	require.NotEmpty(t, configFilePath)
}
