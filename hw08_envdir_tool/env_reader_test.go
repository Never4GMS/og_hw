package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	expectedEnv := Environment{
		"BAR":   {Value: "bar", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: false},
		"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": {Value: "\"hello\"", NeedRemove: false},
		"UNSET": {Value: "", NeedRemove: true},
	}

	t.Run("Read default", func(t *testing.T) {
		env, err := ReadDir("testdata/env")

		require.Nil(t, err)
		for key, expected := range expectedEnv {
			require.Equal(t, expected, env[key], "%s should be %v", key, expected)
		}
	})
}
