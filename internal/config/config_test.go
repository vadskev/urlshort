package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMustLoad(t *testing.T) {
	cgf := MustLoad()
	require.NotEmpty(t, cgf)
}
