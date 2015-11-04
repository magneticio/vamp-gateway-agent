package main

import (
	"io"
	"os"
	go_logging "github.com/op/go-logging"
)

func CreateLogger() *go_logging.Logger {
	var logger = go_logging.MustGetLogger("vamp-gateway-agent")
	var backend = go_logging.NewLogBackend(io.Writer(os.Stdout), "", 0)
	backendFormatter := go_logging.NewBackendFormatter(backend, go_logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortpkg:.4s} %{level:.4s} ==> %{message} %{color:reset}",
	))
	go_logging.SetBackend(backendFormatter)
	return logger
}
