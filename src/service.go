package main

import (
	"errors"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
)

type service struct{}

func (s *service) Execute(_ []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	var tick = time.Tick(config.Interval * time.Minute)
	logger.Info("Сервис запущен.")
	status <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

loop:
	for {
		select {
		case <-tick:
			for _, paths := range config.Paths {
				if _, err := os.Stat(paths.Import); err == nil {
					var data []byte

					data, err = os.ReadFile(paths.Import)
					if err != nil {
						logger.Error("При чтении ATOL-отчета была получена ошибка!", zap.Error(err), zap.String("path", paths.Import))
						continue
					}

					var text = string(data)
					var rows = strings.Split(text, "\n")
					for i, row := range rows {
						var cells = strings.Split(row, ";")
						if len(cells) < 16 {
							continue
						}
						if cells[7] == "" {
							continue
						}
						var temp = cells[7]
						cells[7] = cells[15]
						cells[15] = temp
						rows[i] = strings.Join(cells, ";")
					}
					text = strings.Join(rows, "\n")
					data = []byte(text)

					err = os.WriteFile(paths.Export, data, 0644)
					if err != nil {
						logger.Error("При записи ASTOR-отчета была получена ошибка!", zap.Error(err), zap.String("path", paths.Export))
						continue
					}

					err = os.Remove(paths.Import)
					if err != nil {
						logger.Error("При удалении ATOL-отчета была получена ошибка!", zap.Error(err), zap.String("path", paths.Import))
						continue
					}
					logger.Error("Отчет конвертирован и оригинал удален",
						zap.String("original", paths.Import), zap.String("converted", paths.Export))
				} else if errors.Is(err, os.ErrNotExist) {
					logger.Error("Конвертация пропущена: отсутствует файл", zap.String("path", paths.Import))
				} else {
					logger.Error("Конвертация пропущена: получена ошибка", zap.Error(err), zap.String("path", paths.Import))
				}
			}
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				logger.Info("Остановка службы!")
				break loop
			default:
				logger.Error("Получен некорректный запрос службе...", zap.Uint32("request", uint32(c.Cmd)))
			}
		}
	}

	status <- svc.Status{State: svc.StopPending}
	return false, 0
}
