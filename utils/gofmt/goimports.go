package gofmt

import (
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GoFormatDirectory run gofmt in directory `dir`, every *.go file will be reformatted.
func GoFormatDirectory(dir string) error {
	err := filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(fpath, ".go") && !info.IsDir() {
			err := GoFormatFile(fpath)
			if err != nil {
				log.Printf("Warn: style file:%s error:%v", fpath, err)
			}
		}
		return nil
	})
	return err
}

// GoFormatDirectory reformat file `fpath`, `fpath` must be a valid *.go file.
func GoFormatFile(fpath string) error {

	in, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	out, err := format.Source(in)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	err = ioutil.WriteFile(fpath, out, 0644)
	if err != nil {
		return err
	}

	return nil
}
