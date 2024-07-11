package gimbap

import (
	"github.com/jhseong7/gimbap/app"
)

// Public entry point to create a Gimbap application.
func CreateApp(option app.AppOption) *app.GimbapApp {
	return app.CreateApp(option)
}

// Public to retrive the injected provider.
func GetProvider[T interface{}](a app.GimbapApp, prov T) (ret T, err error) {
	return app.GetProvider(a, prov)
}
