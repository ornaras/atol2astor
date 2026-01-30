package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"go.uber.org/zap"
)

var (
	openDialog cfd.OpenFileDialog
	saveDialog cfd.SaveFileDialog

	openDialogArgs = cfd.DialogConfig{
		Title:       "Выберите ATOL-отчет",
		FileFilters: []cfd.FileFilter{{"Текстовый файл", "*.txt"}},
		Role:        "ImportFrontol6",
	}
	saveDialogArgs = cfd.DialogConfig{
		Title:       "Выберите путь для сохранения ASTOR-отчет",
		FileFilters: []cfd.FileFilter{{"Текстовый файл", "*.txt"}},
		Role:        "ExportAstor1.3",
	}
)

func showConfigurator() {
	var err error
	var interval int64
	var numPaths int

	if openDialog, err = cfd.NewOpenFileDialog(openDialogArgs); err != nil {
		logger.Error("Не удалось создать OpenFileDialog", zap.Error(err))
		fmt.Println("Не удалось проинициализировать конфигуратор")
		return
	}
	defer func() {
		if err = openDialog.Release(); err != nil {
			logger.Error("Не удалось освободить OpenFileDialog", zap.Error(err))
		}
	}()

	if saveDialog, err = cfd.NewSaveFileDialog(saveDialogArgs); err != nil {
		logger.Error("Не удалось создать SaveFileDialog", zap.Error(err))
		fmt.Println("Не удалось проинициализировать конфигуратор")
		return
	}
	defer func() {
		if err = saveDialog.Release(); err != nil {
			logger.Error("Не удалось освободить SaveFileDialog", zap.Error(err))
		}
	}()

	if interval, err = inputInt64("Введите длительность интервала между проверками файлов (в минутах): "); err != nil {
		logger.Error("Не удалось установить длительность интервала между проверками файлов", zap.Error(err))
		fmt.Println("Не удалось установить длительность интервала между проверками файлов")
		return
	}

	if numPaths, err = inputInt("Введите количество пар путей импорт-экспорт: "); err != nil {
		logger.Error("Не удалось определить количество пар путей импорт-экспорт", zap.Error(err))
		fmt.Println("Не удалось определить количество пар путей импорт-экспорт")
		return
	}

	config = &configuration{
		Interval: time.Duration(interval),
		Paths:    make([]*reportPath, numPaths),
	}

	for i := 0; i < numPaths; i++ {
		im, ex, ok := selectPaths()
		if !ok {
			return
		}

		config.Paths[i] = &reportPath{im, ex}
	}

	config.save()
}

func selectPaths() (pathImport, pathExport string, isOk bool) {
	var err error

	if err = openDialog.Show(); err != nil {
		logger.Error("Не удалось показать OpenFileDialog", zap.Error(err))
		fmt.Println("Не удалось открыть выбор пути к ATOL-отчету")
		return
	}
	if pathImport, err = openDialog.GetResult(); err != nil {
		logger.Error("Не удалось определить путь к файлу ATOL-формат", zap.Error(err))
		fmt.Println("Не удалось определить путь к файлу ATOL-формат")
		return
	}

	if err = saveDialog.Show(); err != nil {
		logger.Error("Не удалось показать SaveFileDialog", zap.Error(err))
		fmt.Println("Не удалось открыть выбор пути для ASTOR-отчета")
		return
	}
	if pathExport, err = saveDialog.GetResult(); err != nil {
		logger.Error("Не удалось определить путь для ASTOR-отчета", zap.Error(err))
		fmt.Println("Не удалось определить путь для ASTOR-отчета")
		return
	}
	isOk = true
	return
}

func inputInt64(text string) (int64, error) {
	var err error
	var input string
	var res int64

	fmt.Print(text)
	if _, err = fmt.Scanln(&input); err != nil {
		return 0, err
	}
	if res, err = strconv.ParseInt(input, 10, 64); err != nil {
		return 0, err
	}
	return res, nil
}

func inputInt(text string) (int, error) {
	var err error
	var input string
	var res int

	fmt.Print(text)
	if _, err = fmt.Scanln(&input); err != nil {
		return 0, err
	}
	if res, err = strconv.Atoi(input); err != nil {
		return 0, err
	}
	return res, nil
}
