package wal

import (
	"bufio"
	"fmt"
	"os"
)

type WAL struct {
	file *os.File
}

func NewWAL(f *os.File) *WAL {
	return &WAL{
		file: f,
	}
}

func (w *WAL) Write(cmd string) error {
	_, err := fmt.Fprintf(w.file, "%s\n", cmd)
	if err != nil {
		return err
	}

	return nil
}

func (w *WAL) Read() ([]string, error) {
	var res []string

	scanner := bufio.NewScanner(w.file)
	for scanner.Scan() {
		line := scanner.Text()
		res = append(res, line)
	}

	return res, nil
}

func (w *WAL) Clear() error {
	if err := w.file.Truncate(0); err != nil {
		return err
	}

	if _, err := w.file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}
