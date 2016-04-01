package logger

import (
	"os"
	"github.com/stvp/go-toml-config"
)

var (
	logPath string
	logFile string
	enableConsole bool
	level LEVEL
)

func InitLogger(path string) {
	loadConfig(path)

	_, err := os.Stat(logPath)
	if !os.IsExist(err) {
		err = os.MkdirAll(logPath, os.ModePerm)
	}

	SetConsole(enableConsole)
	SetRollingDaily(logPath, logFile)
	SetLevel(level)

	if err != nil {
		Fatal("Init log path error", err)
	} else {
		Info("Log path initialised.")
	}
}

func loadConfig(path string) {
	logConfig := config.NewConfigSet("logConfig", config.ExitOnError)
	logConfig.StringVar(&logPath, "log_path", "./logs/")
	logConfig.StringVar(&logFile, "log_file", "app.log")
	logConfig.BoolVar(&enableConsole, "console", true)
	var tempLeveStr string
	logConfig.StringVar(&tempLeveStr, "level", "INFO")
	switch tempLeveStr {
	case "ALL":
		level = ALL
	case "DEBUG":
		level = DEBUG
	case "INFO":
		level = INFO
	case "WARN":
		level = WARN
	case "ERROR":
		level = ERROR
	case "FATAL":
		level = FATAL
	case "OFF":
		level = OFF
	}

	err := logConfig.Parse(path)
	if err != nil {
		Warnf("load dbconfig error, %v", err)
	} else {
		Info("loaded dbconfig")
	}
}