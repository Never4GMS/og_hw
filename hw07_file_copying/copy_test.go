package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	src, err := os.Open("testdata/input.txt")
	if err != nil {
		t.Error(err.Error())
	}
	defer src.Close()
	removeTmp := func(f *os.File) {
		n := f.Name()
		f.Close()
		os.Remove(n)
	}
	compareWith := func(t *testing.T, expectedPath, actualPath string) {
		expected, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Error(err.Error())
		}
		actual, err := os.ReadFile(actualPath)
		if err != nil {
			t.Error(err.Error())
		}

		require.Equal(t, len(expected), len(actual))
		require.Equal(t, expected, actual)
	}

	t.Run("Copy full", func(t *testing.T) {
		dst, err := os.CreateTemp("/tmp", "go_hw_07_copy_*.txt")
		if err != nil {
			t.Error(err.Error())
		}
		defer removeTmp(dst)

		err = Copy(dst, src, 0, 10000)
		require.Nil(t, err)
		compareWith(t, "testdata/out_offset0_limit0.txt", dst.Name())
	})

	t.Run("Copy with limit 10", func(t *testing.T) {
		dst, err := os.CreateTemp("/tmp", "go_hw_07_copy_*.txt")
		if err != nil {
			t.Error(err.Error())
		}
		defer removeTmp(dst)

		err = Copy(dst, src, 0, 10)
		require.Nil(t, err)
		compareWith(t, "testdata/out_offset0_limit10.txt", dst.Name())
	})
}
