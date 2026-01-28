package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func installService() {
	var appPath = path.Join(pathWorkdir, "atol2astor.exe")
	var err error
	var originPath string

	if originPath, err = getExePath(); err != nil {
		logger.Error("Не удалось получить путь к исполнительному файлу.", zap.Error(err))
		fmt.Println("Не удалось установить службу.")
		return
	}
	if err = copyFile(originPath, appPath); err != nil {
		logger.Error("Не удалось скопировать исполнительный файл.", zap.Error(err))
		fmt.Println("Не удалось установить службу.")
		return
	}

	var m *mgr.Mgr
	m, err = mgr.Connect()
	if err != nil {
		logger.Error("Не удалось подключиться к менеджеру служб.", zap.Error(err))
		fmt.Println("Не удалось установить службу.")
		return
	}
	defer func() {
		if err = m.Disconnect(); err != nil {
			logger.Error("Не удалось отключиться от менеджера служб.", zap.Error(err))
		}
	}()

	var s *mgr.Service
	if s, err = m.OpenService(nameService); err == nil {
		logger.Error("Не удалось подключиться к службе.", zap.Error(err))
		fmt.Println("Не удалось установить службу.")
		return
	}
	defer func() {
		if err = s.Close(); err != nil {
			logger.Error("Не удалось отключиться от службы.", zap.Error(err))
		}
	}()

	var confService = mgr.Config{
		StartType:   mgr.StartAutomatic,
		DisplayName: "Конвертер ATOL2ASTOR",
	}
	if s, err = m.CreateService(nameService, appPath, confService, "-s"); err != nil {
		logger.Error("Не удалось установить службу.", zap.Error(err))
		fmt.Println("Не удалось установить службу.")
		return
	}

	var events uint32 = eventlog.Error | eventlog.Warning | eventlog.Info
	if err = eventlog.InstallAsEventCreate(nameService, events); err != nil {
		logger.Error("Не удалось установить событие установки.", zap.Error(err))
		fmt.Println("Не удалось установить службу.")
		if err = s.Delete(); err != nil {
			logger.Error("Не удалось удалить службу.", zap.Error(err))
		}
		return
	}

	if err = s.Start(); err != nil {
		logger.Error("Не удалось запустить службу.", zap.Error(err))
		fmt.Println("Не удалось запустить службу.")
		return
	}
}

func uninstallService() {
	var m *mgr.Mgr
	var err error

	if m, err = mgr.Connect(); err != nil {
		logger.Error("Не удалось подключиться к менеджеру служб.", zap.Error(err))
		fmt.Println("Не удалось удалить службу.")
		return
	}
	defer func() {
		if err = m.Disconnect(); err != nil {
			logger.Error("Не удалось отключиться от менеджера служб.", zap.Error(err))
		}
	}()

	var s *mgr.Service
	if s, err = m.OpenService(nameService); err != nil {
		logger.Error("Не удалось подключиться к службе.", zap.Error(err))
		fmt.Println("Не удалось удалить службу.")
		return
	}
	defer func() {
		if err = s.Close(); err != nil {
			logger.Error("Не удалось отключиться от службы.", zap.Error(err))
		}
	}()

	if _, err = s.Control(svc.Stop); err != nil {
		logger.Error("Не удалось остановить службу.", zap.Error(err))
		fmt.Println("Не удалось остановить службу.")
		return
	}

	if err = s.Delete(); err != nil {
		logger.Error("Не удалось удалить службу.", zap.Error(err))
		fmt.Println("Не удалось удалить службу.")
		return
	}

	if err = eventlog.Remove(nameService); err != nil {
		logger.Error("Не удалось удалить событие об установке службы.", zap.Error(err))
		fmt.Println("Не удалось удалить службу.")
		return
	}
}

func openConfiguration() {
	go func() {
		var err error
		var notepadPath = path.Join(os.Getenv("WINDIR"), "notepad.exe")
		var filePath = path.Join(pathWorkdir, confFileName)

		if err = exec.Command(notepadPath, filePath).Run(); err != nil {
			logger.Error("Не удалось открыть файл конфигурации.", zap.Error(err))
			fmt.Println("Не удалось открыть файл конфигурации.")
		}
	}()
	time.Sleep(time.Second / 10)
}
