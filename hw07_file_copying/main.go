package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	src, err := os.OpenFile(from, os.O_RDONLY, 0o444)
	if err != nil {
		fmt.Printf("copy: can't open from file: %v", err.Error())
		return
	}
	defer src.Close()

	dst, err := os.Create(to)
	if err != nil {
		fmt.Printf("copy: can't create to file: %v", err.Error())
		return
	}
	defer dst.Close()

	s, err := src.Stat()
	if err != nil {
		fmt.Printf("copy: can't get from file info: %v", err.Error())
		return
	}

	if limit == 0 || (offset+limit) > s.Size() {
		limit = s.Size() - offset
	}

	bar := progressbar.DefaultBytes(limit, "copying...")
	defer bar.Finish()
	if err := Copy(io.MultiWriter(dst, bar), src, offset, limit); err != nil {
		fmt.Printf("failed to copy: %v", err)
	}
}
