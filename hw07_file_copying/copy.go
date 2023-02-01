package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrLimitExceedsFileSize  = errors.New("limits smaller file size sub offset")
	ErrInvalidPath           = errors.New("invalid path")
	ErrFileIsNilPointer      = errors.New("use nil pointer as file")
	ErrTmpCreate             = errors.New("cannot create tmp file")
)

const BUF_SIZE = 64 //kB

type FileCop struct {
	file *os.File
	out  *os.File
	info os.FileInfo
}

func NewFileCop(fromPath string) (*FileCop, error) {
	src, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Print(time.Now(), err)
		return nil, ErrInvalidPath
	}

	info, err := src.Stat()
	if err != nil {
		log.Print(time.Now(), err)
		return nil, ErrUnsupportedFile
	}

	if info.IsDir() {
		return nil, ErrUnsupportedFile
	}

	out, err := os.CreateTemp("tmp_copy", "*")
	if err != nil {
		log.Print(time.Now(), err)
		return nil, ErrInvalidPath
	}

	return &FileCop{
		file: src,
		out:  out,
		info: info,
	}, nil
}

func (cf *FileCop) checkOffsetAndLimit(offset, limit int64) error {
	if cf.info.Size() <= offset || cf.info.Size() <= 0 {
		return ErrOffsetExceedsFileSize
	}

	if cf.info.Size()-offset < limit {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

func (cf *FileCop) Copy(offset, limit int64, progress chan uint8) {
	if limit < 0 {
		return
	}
	step := (cf.info.Size() / (BUF_SIZE * 1024)) + 1
	buf := make([]byte, BUF_SIZE*1024)

	for currentStep := int64(0); currentStep < step; currentStep++ {
		n, err := cf.file.ReadAt(buf, offset)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		if n > 0 {
			x, err := cf.out.Write(buf)
			if err != nil {
				panic(err)
			}

			if n != x {
				log.Printf("invalid step - scan %d, write %d", n, x)
			}
			progress <- uint8((currentStep / step) * 100)
		}
	}

	close(progress)
}

func (cf *FileCop) ReadInBuf(offset, limit int64) (chan uint8, error) {
	progress := make(chan uint8, 100)

	if err := cf.checkOffsetAndLimit(offset, limit); err != nil {
		return nil, err
	}

	go cf.Copy(offset, limit, progress)

	return progress, nil
}

func (cf *FileCop) Close() {
	cf.file.Close()
	cf.out.Close()
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrInvalidPath
	}

	cf, err := NewFileCop(fromPath)
	if err != nil {
		return err
	}

	progress, err := cf.ReadInBuf(offset, limit)

	for step := range progress {
		fmt.Printf("\rprogress - %d %%", step)
	}

	cf.Close()

	if err != nil {
		return err
	}

	return nil
}
