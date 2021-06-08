package session

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type MessageWriter interface {
	Init() error
	Write([]byte) error
}

type FilenameGenerator func() string

func NewMessageWriter(baseDir string, sessionID string, gen FilenameGenerator) MessageWriter {
	if baseDir == "" {
		return &messageDummyWriter{}
	} else {
		return &messageFileWriter{
			dir:               filepath.Join(baseDir, sessionID),
			filenameGenerator: gen,
		}
	}
}

type messageFileWriter struct {
	dir               string
	filenameGenerator FilenameGenerator
}

func (w *messageFileWriter) Init() error {
	return os.MkdirAll(w.dir, 0700)
}

func (w *messageFileWriter) Write(message []byte) error {
	filename := filepath.Join(w.dir, w.filenameGenerator())
	return ioutil.WriteFile(filename, message, 0660)
}

type messageDummyWriter struct{}

func (w *messageDummyWriter) Init() error {
	return nil
}

func (w *messageDummyWriter) Write([]byte) error {
	return nil
}
