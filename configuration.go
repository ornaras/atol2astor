package main

import (
	"bytes"
	"encoding/xml"
	"os"
	"time"

	"go.uber.org/zap"
)

var config *configuration

const confFileName = "config.xml"

type configuration struct {
	IntervalComment    string        `xml:",comment"`
	Interval           time.Duration `xml:"interval"`
	ImportPathsComment string        `xml:",comment"`
	ExportPathsComment string        `xml:",comment"`
	Warning1           string        `xml:",comment"`
	Warning2           string        `xml:",comment"`
	ImportPaths        []*importPath `xml:"imports>import"`
}

type importPath struct {
	Path string `xml:"path,attr"`
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
	var buff = new(bytes.Buffer)
	var encoder = xml.NewEncoder(buff)
	defer encoder.Close()
	encoder.Indent("", "    ")
	err = encoder.Encode(conf)
	if err != nil {
		logger.Error("При переводе конфигурации в XML произошла ошибка!", zap.Error(err))
		panic(err)
	}
	var b = bytes.Replace(buff.Bytes(), []byte("></import>"), []byte("/>"), -1)

	if _, err = file.Write(b); err != nil {
		logger.Error("При сохранении конфигурации произошла ошибка!", zap.Error(err))
		panic(err)
	}
}

func createDefaultConfiguration() *configuration {
	return &configuration{
		IntervalComment:    "Тег 'interval' устанавливает интервал между проверками файлов в минутах",
		Interval:           5,
		ImportPathsComment: "В теге 'imports' хранятся пути к файлам, которые будут конвертироваться в формат ASTOR",
		ExportPathsComment: "Конвертированный файл сохраняется в той же директории, в которой находится оригинальный файл, с имененем 'export.txt'",
		Warning1:           "Внимание! Не рекомендуется использовать несколько раз одну и ту же директорию.",
		Warning2:           "Внимание! Оригинальный файл должен иметь имя отличное от 'export.txt'.",
		ImportPaths: []*importPath{
			{"C:\\atol1\\import.txt"},
			{"C:\\atol2\\import.txt"},
		},
	}
}
