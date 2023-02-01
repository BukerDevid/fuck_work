package main

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidPath           = errors.New("invalid path")
	ErrTmpCreate             = errors.New("cannot create tmp file")
)

type OffsetReader struct {
	src io.Reader
}

func prepFile(fromPath string) (src *os.File, out *os.File, err error) {
	src, err = os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Print(time.Now(), err)
		err = ErrInvalidPath
		return
	}

	out, err = os.CreateTemp("tmp_copy", "*")
	if err != nil {
		log.Print(time.Now(), err)
		err = ErrInvalidPath
		return
	}

	return
}

func checkSize(stat func() (fs.FileInfo, error), offset int64) error {
	fileInfo, err := stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrInvalidPath
	}

	src, out, err := prepFile(fromPath)
	if err != nil {
		return err
	}

	defer func() {
		os.Remove(out.Name())
	}()

	if err := checkSize(src.Stat, offset); err != nil {
		return err
	}

	if limit > 0 {
		_, err = io.CopyN(out, src, limit)
	} else {
		_, err = io.Copy(out, src)
	}

	if err != nil {
		return err
	}

	return nil
}
