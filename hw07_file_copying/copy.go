package main

import (
	"errors"
	"io"
)

const defaultBufferLength int64 = 1024

func Copy2(dst io.Writer, src io.ReadSeeker, offset, limit int64) error {
	if offset > 0 {
		if _, err := src.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}
	_, err := io.CopyN(dst, io.Reader(src), limit)
	return err
}

func Copy(dst io.WriterAt, src io.ReaderAt, offset, limit int64) error {
	buffer := make([]byte, min(limit, defaultBufferLength))
	var writeOffset int64
	readNext := true
	for readNext {
		n, err := src.ReadAt(buffer, offset)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			readNext = false
		}
		offset += int64(n)
		if offset >= limit {
			offset = limit
			readNext = false
		}

		rb := buffer
		if n < len(buffer) {
			rb = buffer[:n]
		}
		n, err = dst.WriteAt(rb, writeOffset)
		if err != nil {
			return err
		}
		writeOffset += int64(n)
	}

	return nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
