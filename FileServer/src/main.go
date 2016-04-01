package main

import (
	"net"
	"github.com/kyugao/go-logger/logger"
	_ "logs"
	"google.golang.org/grpc"
	"external"
	"service"
	_ "db"
	"github.com/stvp/go-toml-config"
	"fmt"
)

const service_config_path = "./conf/service.conf"

var (
	serviceConfig *config.ConfigSet
	servicePort string
)

func init() {
	loadConfig()
}

func loadConfig() {
	serviceConfig = config.NewConfigSet("serviceConfig", config.ExitOnError)
	serviceConfig.StringVar(&servicePort, "port", "6010")
	err := serviceConfig.Parse(service_config_path)
	if err != nil {
		logger.Warnf("load service config error, %v", err)
		return
	} else {
		logger.Info("loaded service config")
	}
}

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", servicePort))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	} else {
		logger.Debugf("listening on port %s", servicePort)
	}
	// 创建grpc实例
	grpcServer := grpc.NewServer()
	// 注册fileService服务
	external.RegisterServiceServer(grpcServer, service.InitServer())
	// 启动服务
	grpcServer.Serve(lis)
}