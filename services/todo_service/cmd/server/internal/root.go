package cmd

import (
	"log"
	"os"

	"todo_service/cmd/server/internal/handlers"
	"todo_service/common"

	"github.com/gofiber/fiber/v2"
	goservice "github.com/onpointvn/libs/go-sdk"
	"github.com/onpointvn/libs/go-sdk/plugin/fiberapp"
	"github.com/onpointvn/libs/go-sdk/plugin/storage/sdkgorm"
	"github.com/onpointvn/libs/go-sdk/plugin/storage/sdkredis"
	"github.com/spf13/cobra"
)

var (
	serviceName = "todo-service"
	version     = "1.0.0"
)

func newService() goservice.Service {
	s := goservice.New(
		goservice.WithName(serviceName),
		goservice.WithVersion(version),
		goservice.WithInitRunnable(fiberapp.New(common.PluginFiber)),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main", common.DBMain)),
		goservice.WithInitRunnable(sdkredis.NewRedisDB("main", common.PluginRedis)),
	)

	return s
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start todo service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		service.MustGet(common.PluginFiber).(interface {
			SetRegisterHdl(app *fiber.App)
		}).SetRegisterHdl(handlers.Router(service))

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
