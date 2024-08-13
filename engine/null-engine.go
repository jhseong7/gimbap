package engine

import (
	"github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"
)

type (
	NullEngine struct {
		IServerEngine

		logger ecl.Logger

		stopFlag chan string
	}
)

func (e *NullEngine) RegisterController(rootPath string, instance controller.IController) {
	e.logger.Warn("NullEngine does not support controller registration. Please check if this is intended.")
}

func (e *NullEngine) Run(option ServerRuntimeOption) {
	e.stopFlag = make(chan string)

	// Wait for the stop signal (blocking)
	<-e.stopFlag

	return
}

func (e *NullEngine) Stop() {
	// Send the stop signal to free the Run() method
	e.stopFlag <- "stop"

	return
}

func NewNullEngine() *NullEngine {
	return &NullEngine{
		logger: ecl.NewLogger(ecl.LoggerOption{
			Name: "NullEngine",
		}),
	}
}
