package tool

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
)

var sourceCodeFiles = []string{
	"application.go",
	"application_callback.go",
	"application_config.go",
	"application_event.go",
	"application_queue_async_call.go",
	"application_queue_async_call_posix.go",
	"application_queue_async_call_windows.go",
	"application_run.go",
}

func LoadSourceFile(path string, filter bool) []byte {
	buf := bytes.Buffer{}
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		isExt := filepath.Ext(info.Name()) == ".md" || filepath.Ext(info.Name()) == ".go"
		if !info.IsDir() && isExt {
			if filter {
				i := sort.SearchStrings(sourceCodeFiles, info.Name())
				if i < len(sourceCodeFiles) && sourceCodeFiles[i] == info.Name() {
					data, _ := ioutil.ReadFile(path)
					buf.WriteString("## " + info.Name() + "\n")
					buf.WriteString("\n```go\n")
					buf.Write(data)
					buf.WriteString("\n```\n")
				}
			} else {
				data, _ := ioutil.ReadFile(path)
				buf.WriteString("## " + info.Name() + "\n")
				buf.WriteString("\n```go\n")
				buf.Write(data)
				buf.WriteString("\n```\n")
			}
		}
		return nil
	})
	return buf.Bytes()
}

func LoadFile(path string) []byte {
	buf := bytes.Buffer{}
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		isExt := filepath.Ext(info.Name()) == ".md" || filepath.Ext(info.Name()) == ".go"
		if !info.IsDir() && isExt {
			data, _ := ioutil.ReadFile(path)
			buf.Write(data)
			buf.WriteString("\n")
		}
		return nil
	})
	return buf.Bytes()
}
