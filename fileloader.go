package main

import (
	"io/ioutil"
	"net/http"
)

type FileLoader struct {
	Path string
}

func (fileLoader FileLoader) LoadFile(writer http.ResponseWriter, request *http.Request) {
	form, err := ioutil.ReadFile(fileLoader.Path)

	if err != nil {
		return
	}

	if _, err := writer.Write(form); err != nil {
		return
	}
}
