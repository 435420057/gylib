package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"handler"
	"github.com/stvp/go-toml-config"
	"gylogger"
	"fmt"
	"runtime"
	"gycache"
)


const service_config_path = "./conf/service.conf"

var (
	serviceConfig *config.ConfigSet
	servicePort string
)

func loadConfig() {
	serviceConfig = config.NewConfigSet("serviceConfig", config.ExitOnError)
	serviceConfig.StringVar(&servicePort, "port", "8080")
	err := serviceConfig.Parse(service_config_path)
	if err != nil {
		logger.Warnf("load service config error, %v", err)
		return
	} else {
		logger.Info("loaded service config")
	}
}

func init() {

	logger.InitLogger("./conf/logger.conf")
	cache.InitCache("./conf/cache.conf")

	loadConfig()
	r := mux.NewRouter()
	r.HandleFunc("/f/{action}", handler.FormHandler).Methods("POST")
	r.HandleFunc("/file/upload", handler.UploadHandler).Methods("POST")
	r.HandleFunc("/file/download/{fileId}", handler.DownloadHandler)
	http.Handle("/", r)
}

func main() {

	logger.Infof("====== Start agent service node @ %s ======", servicePort)

	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1 << 20)
			runtime.Stack(buf, true)
			logger.Debug(buf)
		}
	}()

	port := fmt.Sprintf(":%s", servicePort)
	http.ListenAndServe(port, nil)
}
