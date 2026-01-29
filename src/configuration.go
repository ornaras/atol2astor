package main

import (
	"bytes"
	"encoding/xml"
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	config     *configuration
	pathConfig string
)

const confFileName = "config.xml"

type configuration struct {
	IntervalComment    string        `xml:",comment"`
	Interval           time.Duration `xml:"interval"`
	ImportPathsComment string        `xml:",comment"`
	ExportPathsComment string        `xml:",comment"`
	Warning1           string        `xml:",comment"`
	Paths              []*reportPath `xml:"reports>report"`
}

type reportPath struct {
	Import string `xml:"import,attr"`
	Export string `xml:"export,attr"`
}

func (conf *configuration) save() {
	var err error
	var file *os.File

	file, err = os.OpenFile(pathConfig, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		logger.Error("При создании файла конфигурации произошла ошибка!", zap.Error(err))
		panic(err)
	}
	defer file.Close()
	var buff = new(bytes.Buffer)
	var encoder = xml.NewEncoder(buff)
	defer encoder.Close()
	encoder.Indent("", "    ")
	err = encoder.Encode(conf)
	if err != nil {
		logger.Error("При переводе конфигурации в XML произошла ошибка!", zap.Error(err))
		panic(err)
	}
	var b = bytes.Replace(buff.Bytes(), []byte("></report>"), []byte("/>"), -1)

	if _, err = file.Write(b); err != nil {
		logger.Error("При сохранении конфигурации произошла ошибка!", zap.Error(err))
		panic(err)
	}
}

func createDefaultConfiguration() *configuration {
	return &configuration{
		IntervalComment:    "Тег 'interval' устанавливает интервал между проверками файлов в минутах",
		Interval:           5,
		ImportPathsComment: "В теге 'reports' хранятся пути к файлам конвертации",
		Warning1:           "Внимание! Не рекомендуется использовать несколько раз одну и ту же директорию.",
		Paths: []*reportPath{
			{"C:\\atol1\\import.txt", "C:\\atol1\\export.txt"},
			{"C:\\atol2\\import.txt", "C:\\atol2\\export.txt"},
		},
	}
}

func loadConfiguration() {
	var err error
	var file *os.File

	file, err = os.Open(pathConfig)
	if err != nil {
		logger.Error("При открытии файла конфигурации произошла ошибка!", zap.Error(err))
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("При закрытии файла конфигурации произошла ошибка!", zap.Error(err))
		}
	}()

	err = xml.NewDecoder(file).Decode(&config)
	if err != nil {
		logger.Error("При переводе XML в конфигурацию произошла ошибка!", zap.Error(err))
		panic(err)
	}
}
