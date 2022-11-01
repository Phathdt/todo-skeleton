package fiberapp

import (
	"flag"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/onpointvn/libs/go-sdk/logger"
)

type fiberApp struct {
	prefix string
	port   int
	logger logger.Logger
	app    *fiber.App
}

func New(prefix string) *fiberApp {
	return &fiberApp{prefix: prefix}
}

func (f *fiberApp) GetPrefix() string {
	return f.prefix
}

func (f *fiberApp) Get() interface{} {
	return f
}

func (f *fiberApp) Name() string {
	return f.prefix
}

func (f *fiberApp) InitFlags() {
	flag.IntVar(&f.port, f.GetPrefix()+"-port", 4000, "fiber port")
}

func (f *fiberApp) Configure() error {
	f.logger = logger.GetCurrent().GetLogger(f.Name())

	go func() {
		if err := f.app.Listen(fmt.Sprintf(":%d", f.port)); err != nil {
			f.logger.Error(err)
		}
	}()

	return nil
}

func (f *fiberApp) Run() error {
	go func() {
		time.Sleep(time.Second * 5)

		_ = f.Configure()
	}()
	return nil
}

func (f *fiberApp) Stop() <-chan bool {
	if f.app != nil {
		if err := f.app.Shutdown(); err != nil {
			f.logger.Error(err)
		}
	}

	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (f *fiberApp) SetRegisterHdl(app *fiber.App) {
	f.app = app
}
