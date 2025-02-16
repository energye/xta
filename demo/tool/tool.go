package tool

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

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
