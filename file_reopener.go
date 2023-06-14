package ltsvlog

import (
	"io"
	"os"
)

type FileReopener struct {
	file *os.File
	flag int
	perm os.FileMode
	sw   *SwitchableWriter
}

var _ io.Writer = (*FileReopener)(nil)

func NewFileReopener(name string, flag int, perm os.FileMode) (*FileReopener, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return &FileReopener{
		file: file,
		flag: flag,
		perm: perm,
		sw:   NewSwitchableWriter(file),
	}, nil
}

func (h *FileReopener) Write(p []byte) (n int, err error) {
	return h.sw.Write(p)
}

func (h *FileReopener) Reopen() error {
	newFile, err := os.OpenFile(h.file.Name(), h.flag, h.perm)
	if err != nil {
		return err
	}

	h.sw.Switch(newFile)

	if err := h.SyncAndClose(); err != nil {
		return err
	}
	h.file = newFile
	return nil
}

func (h *FileReopener) SyncAndClose() error {
	if err := h.file.Sync(); err != nil {
		return err
	}
	if err := h.file.Close(); err != nil {
		return err
	}
	return nil
}
