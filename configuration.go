package main

import (
	"encoding/xml"
	"os"
	"time"

	"go.uber.org/zap"
)

var config *configuration

const confFileName = "config.xml"

type configuration struct {
	Interval time.Duration `xml:"interval"`
	Paths    []filePath    `xml:"paths>path"`
}

type filePath struct {
	Import string `xml:"import"`
	Export string `xml:"export"`
}

func (conf *configuration) save(path string) {
	var err error
	var file *os.File

	file, err = os.Create(path)
	if err != nil {
		logger.Error("При создании файла конфигурации произошла ошибка!", zap.Error(err))
		panic(err)
	}
	defer file.Close()
	err = xml.NewEncoder(file).Encode(&config)
	if err != nil {
		logger.Error("При переводе конфигурации в XML произошла ошибка!", zap.Error(err))
		panic(err)
	}
}

func createDefaultConfiguration() *configuration {
	return &configuration{
		Interval: 5,
		Paths: []filePath{
			{
				Import: "C:\\atol.txt",
				Export: "C:\\astor.txt",
			},
			{
				Import: "C:\\atol2.txt",
				Export: "C:\\astor2.txt",
			},
		},
	}
}
