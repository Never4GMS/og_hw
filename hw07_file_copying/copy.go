package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrUnsupportedDestFile   = errors.New("unsupported dest file")
	ErrSameFile              = errors.New("from and to path are the same")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var err error
	if !filepath.IsAbs(fromPath) {
		if fromPath, err = filepath.Abs(fromPath); err != nil {
			return err
		}
	}
	if !filepath.IsAbs(toPath) {
		if toPath, err = filepath.Abs(toPath); err != nil {
			return err
		}
	}

	if fromPath == toPath {
		return ErrSameFile
	}

	src, err := os.OpenFile(fromPath, os.O_RDONLY, 0o444)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer src.Close()

	s, err := src.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	if offset > s.Size() {
		return ErrOffsetExceedsFileSize
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedDestFile
	}
	defer dst.Close()

	if limit == 0 || (offset+limit) > s.Size() {
		limit = s.Size() - offset
	}

	bar := progressbar.DefaultBytes(limit, "copying...")
	defer bar.Finish()
	return copyFromReaderToWriter(io.MultiWriter(dst, bar), src, offset, limit)
}

func copyFromReaderToWriter(dst io.Writer, src io.ReadSeeker, offset, limit int64) error {
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
