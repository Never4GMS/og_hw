package main

import (
	"errors"
	"io"
)

func Copy(dst io.Writer, src io.ReadSeeker, offset, limit int64) error {
	if offset > 0 {
		if _, err := src.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}
	_, err := io.CopyN(dst, io.Reader(src), limit)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
