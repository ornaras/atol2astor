package main

import (
	"encoding/json"
	"os"
	"time"

	"go.uber.org/zap"
)

var config *configuration

type configuration struct {
	Interval time.Duration `json:"interval"`
	Paths    []filePath    `json:"paths"`
}

type filePath struct {
	Import string `json:"import"`
	Export string `json:"export"`
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
	err = json.NewEncoder(file).Encode(&config)
	if err != nil {
		logger.Error("При переводе конфигурации в JSON произошла ошибка!", zap.Error(err))
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
