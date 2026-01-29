package main

import (
	"flag"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initWorkdir() {
	var programdata = os.Getenv("PROGRAMDATA")
	pathWorkdir = path.Join(programdata, "atol2astor")
	if _, err := os.Stat(pathWorkdir); os.IsNotExist(err) {
		err = os.Mkdir(pathWorkdir, 0777)
		if err != nil {
			panic(err)
		}
	}
	pathConfig = path.Join(pathWorkdir, confFileName)
}

func initLogger() {
	var _path = path.Join(pathWorkdir, "log.txt")
	var err error
	var file *os.File
	file, err = os.OpenFile(_path, os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.LineEnding = "\r\n"
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)

	logger = zap.New(core)
	logger.Info("Инициализация atol2astor " + version)
}

func initFlags() {
	flag.BoolVar(&isService, "s", false, "Запуск в режиме сервиса")
	flag.BoolVar(&isDebug, "d", false, "Запуск в режиме отладки")
}

func initConfig() {
	var err error
	if _, err = os.Stat(pathConfig); os.IsNotExist(err) {
		config = createDefaultConfiguration()
		config.save()
	} else {
		loadConfiguration()
	}
}
