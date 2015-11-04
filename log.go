package main

import (
	"io"
	"os"
	gologging "github.com/op/go-logging"
)


func CreateLogger() *gologging.Logger {
	var logger = gologging.MustGetLogger("vamp-gateway-agent")
	var backend = gologging.NewLogBackend(io.Writer(os.Stdout), "", 0)
	backendFormatter := gologging.NewBackendFormatter(backend, gologging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{shortpkg:.4s} %{level:.4s} ==> %{message} %{color:reset}",
	))
	gologging.SetBackend(backendFormatter)
	return logger
}
