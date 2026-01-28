package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
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
		var m *mgr.Mgr
		var input string
		for {
			fmt.Println("Возможные действия:")
			fmt.Println("1) Установка службы")
			fmt.Println("2) Удаление службы")
			fmt.Println("3) Открыть конфигурацию")
			fmt.Println()
			fmt.Print("Номер действия: ")
			fmt.Scanln(&input)
			switch input {
			case "1":
				var _path string
				_path, err = exePath()
				if err != nil {
					logger.Error("Не удалось получить путь к исполнительному файлу.", zap.Error(err))
					fmt.Println("Не удалось установить службу.")
					m.Disconnect()
					continue
				}
				m, err = mgr.Connect()
				if err != nil {
					logger.Error("Не удалось подключиться к менеджеру служб.", zap.Error(err))
					fmt.Println("Не удалось установить службу.")
					m.Disconnect()
					continue
				}
				s, err := m.OpenService(nameService)
				if err == nil {
					s.Close()
					logger.Error("Не удалось подключиться к службе.", zap.Error(err))
					fmt.Println("Не удалось установить службу.")
					s.Close()
					m.Disconnect()
					continue
				}
				s, err = m.CreateService(nameService, _path, mgr.Config{StartType: mgr.StartAutomatic}, "-s")
				if err != nil {
					logger.Error("Не удалось установить службу.", zap.Error(err))
					fmt.Println("Не удалось установить службу.")
					s.Close()
					m.Disconnect()
					continue
				}
				err = eventlog.InstallAsEventCreate(nameService, eventlog.Error|eventlog.Warning|eventlog.Info)
				if err != nil {
					s.Delete()
					logger.Error("Не удалось установить событие установки.", zap.Error(err))
					fmt.Println("Не удалось установить службу.")
					s.Close()
					m.Disconnect()
					continue
				}
				err = s.Start()
				if err != nil {
					logger.Error("Не удалось запустить службу.", zap.Error(err))
					fmt.Println("Не удалось запустить службу.")
					s.Close()
					m.Disconnect()
					continue
				}
				s.Close()
				m.Disconnect()
				return
			case "2":
				m, err := mgr.Connect()
				if err != nil {
					logger.Error("Не удалось подключиться к менеджеру служб.", zap.Error(err))
					fmt.Println("Не удалось удалить службу.")
					m.Disconnect()
					continue
				}
				s, err := m.OpenService(nameService)
				if err != nil {
					logger.Error("Не удалось подключиться к службе.", zap.Error(err))
					fmt.Println("Не удалось удалить службу.")
					s.Close()
					m.Disconnect()
					continue
				}
				_, err = s.Control(svc.Stop)
				if err != nil {
					logger.Error("Не удалось остановить службу.", zap.Error(err))
					fmt.Println("Не удалось остановить службу.")
					s.Close()
					m.Disconnect()
					continue
				}
				err = s.Delete()
				if err != nil {
					logger.Error("Не удалось удалить службу.", zap.Error(err))
					fmt.Println("Не удалось удалить службу.")
					s.Close()
					m.Disconnect()
					continue
				}
				err = eventlog.Remove(nameService)
				if err != nil {
					logger.Error("Не удалось удалить событие об установке службы.", zap.Error(err))
					fmt.Println("Не удалось удалить службу.")
					s.Close()
					m.Disconnect()
					continue
				}
				s.Close()
				m.Disconnect()
				return
			case "3":
				var _path = path.Join(os.Getenv("WINDIR"), "notepad.exe")
				err = exec.Command(_path, path.Join(pathWorkdir, "config.json")).Run()
				if err != nil {
					logger.Error("Не удалось открыть файл конфигурации.", zap.Error(err))
					fmt.Println("Не удалось открыть файл конфигурации.")
					continue
				}
				return
			}
		}
	}
}
