package session

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type MessageWriter interface {
	Init() error
	Write([]byte) error
}

type FilenameGenerator func(ext string) string

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
	var ext = "txt"
	if isJSON(message) {
		ext = "json"
	}
	filename := filepath.Join(w.dir, w.filenameGenerator(ext))
	return ioutil.WriteFile(filename, message, 0660)
}

func isJSON(data []byte) bool {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return false
	}

	switch obj.(type) {
	case map[string]interface{}:
		return true
	case []interface{}:
		return true
	default:
		return false
	}
}

type messageDummyWriter struct{}

func (w *messageDummyWriter) Init() error {
	return nil
}

func (w *messageDummyWriter) Write([]byte) error {
	return nil
}
