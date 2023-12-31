package main

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/todo-service/internal/config"
	"github.com/todo-lists-app/todo-service/internal/service"
)

var (
	BuildVersion = "dev"
	BuildHash    = "none"
	ServiceName  = "todo-service"
)

func main() {
	logs.Local().Info(fmt.Sprintf("Starting %s", ServiceName))
	logs.Local().Info(fmt.Sprintf("Version: %s, Hash: %s", BuildVersion, BuildHash))

	cfg, err := config.Build()
	if err != nil {
		_ = logs.Errorf("config: %v", err)
		return
	}

	if err := service.NewService(cfg).Start(); err != nil {
		_ = logs.Errorf("service: %v", err)
		return
	}
}
