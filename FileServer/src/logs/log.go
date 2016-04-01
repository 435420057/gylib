package logs
import (
	"github.com/kyugao/go-logger/logger"
	"os"
)

const (
	LOG_PATH = "./logs/"
	LOG_FILE = "file.log"
)

func init() {
	_, err := os.Stat(LOG_PATH)
	if !os.IsExist(err) {
		err = os.MkdirAll(LOG_PATH, os.ModePerm)
	}

	logger.SetConsole(true)
	logger.SetRollingDaily(LOG_PATH, LOG_FILE)
	logger.SetLevel(logger.ALL)

	if err != nil {
		logger.Fatal("Init log path error", err)
	} else {
		logger.Info("Log path initialised.")
	}
}