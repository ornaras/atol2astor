package main

import (
	"flag"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

const nameService = "atol2astor"

var (
	logger      *zap.Logger
	pathWorkdir string

	isService bool
	isDebug   bool
)

func init() {
	initWorkdir()
	initLogger()
	initFlags()
	initConfig()
}

func main() {
	flag.Parse()
	var err error
	if isService {
		err = svc.Run(nameService, &service{})
		if err != nil {
			logger.Error("Во время отладки произошла ошибка!", zap.Error(err))
		}
	} else if isDebug {
		err = debug.Run(nameService, &service{})
		if err != nil {
			logger.Error("Во время работы сервиса произошла ошибка!", zap.Error(err))
		}
	} else {
		var input string
		var err error
		for {
			fmt.Println("Возможные действия:")
			fmt.Println("1) Установка службы")
			fmt.Println("2) Удаление службы")
			fmt.Println("3) Открыть конфигурацию")
			fmt.Println()
			fmt.Print("Номер действия: ")
			if _, err = fmt.Scanln(&input); err != nil {
				logger.Error("Не удалось получить текст ввода.", zap.Error(err))
				fmt.Println("Не удалось получить текст ввода.")
				continue
			}
			switch input {
			case "1":
				installService()
				return
			case "2":
				uninstallService()
				return
			case "3":
				openConfiguration()
				return
			}
		}
	}
}
