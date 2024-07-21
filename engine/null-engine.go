package engine

import (
	logger "github.com/jhseong7/ecl"
	"github.com/jhseong7/gimbap/controller"
)

type (
	NullEngine struct {
		IServerEngine

		logger logger.Logger

		stopFlag chan string
	}
)

func (e *NullEngine) RegisterController(rootPath string, instance controller.IController) {
	e.logger.Warn("NullEngine does not support controller registration. Please check if this is intended.")
}

func (e *NullEngine) Run(port int) {
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
		logger: logger.NewLogger(logger.LoggerOption{
			Name: "NullEngine",
		}),
	}
}
