package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	expectedEnv := Environment{
		"HELLO": {Value: "hello", NeedRemove: false},
		"WORLD": {Value: "world", NeedRemove: false},
	}

	t.Run("Exec printenv", func(t *testing.T) {
		out := bytes.Buffer{}
		code := RunCmd([]string{"/bin/bash", "testdata/exec_test.sh"}, expectedEnv, func(rp *RunParams) { rp.Stdout = &out })

		require.Equal(t, 0, code)
		require.Equal(t, "hello world\n", out.String())
	})
}
