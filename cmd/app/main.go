package main

import (
	"github.com/jhseong7/ecl"

	"github.com/jhseong7/gimbap"
	"github.com/jhseong7/gimbap/engine"
	"github.com/jhseong7/gimbap/microservice"
	"github.com/jhseong7/gimbap/provider"
	sample "github.com/jhseong7/gimbap/sample/http-engine"
)

type (
	// Sample microservice
	SampleMicroService struct {
		microservice.IMicroService
		logger ecl.Logger
	}

	SampleMicroService2 struct {
		Text string
	}
)

func (s *SampleMicroService) Start() {
	s.logger.Log("Starting SampleMicroService")
}

func (s *SampleMicroService) Stop() {
	s.logger.Log("Stopping SampleMicroService")
}

func NewSampleMicroService(sam *SampleMicroService2) *SampleMicroService {
	return &SampleMicroService{
		logger: ecl.NewLogger(ecl.LoggerOption{
			Name: "SampleMicroService" + sam.Text,
		}),
	}
}

func main() {
	app := gimbap.CreateApp(gimbap.AppOption{
		AppName:      "SampleApp",
		AppModule:    sample.AppModuleFiber,
		ServerEngine: engine.NewFiberHttpEngine(),
	})

	app.UseInjection(
		&SampleMicroService2{
			Text: "Injec!",
		},
	)

	app.AddMicroServices(microservice.MicroServiceProvider{
		Provider: provider.Provider{
			Name:         "SampleMicroService",
			Instantiator: NewSampleMicroService,
		},
	})

	app.Run()
}
