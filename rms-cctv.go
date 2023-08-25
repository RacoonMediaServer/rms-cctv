package main

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-cctv/internal/camera"
	"github.com/RacoonMediaServer/rms-cctv/internal/cctv"
	"github.com/RacoonMediaServer/rms-cctv/internal/config"
	"github.com/RacoonMediaServer/rms-cctv/internal/db"
	"github.com/RacoonMediaServer/rms-cctv/internal/manager"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactions"
	"github.com/RacoonMediaServer/rms-cctv/internal/reactor"
	"github.com/RacoonMediaServer/rms-cctv/internal/service"
	"github.com/RacoonMediaServer/rms-cctv/internal/settings"
	"github.com/RacoonMediaServer/rms-cctv/internal/timeline"
	"github.com/RacoonMediaServer/rms-packages/pkg/pubsub"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

var Version = "v0.0.0"

const serviceName = "rms-cctv"

func main() {
	logger.Infof("%s %s", serviceName, Version)
	defer logger.Info("DONE.")

	useDebug := false

	microService := micro.NewService(
		micro.Name(serviceName),
		micro.Version(Version),
		micro.Flags(
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"debug"},
				Usage:       "debug log level",
				Value:       false,
				Destination: &useDebug,
			},
		),
	)

	microService.Init(
		micro.Action(func(context *cli.Context) error {
			configFile := fmt.Sprintf("/etc/rms/%s.json", serviceName)
			if context.IsSet("config") {
				configFile = context.String("config")
			}
			return config.Load(configFile)
		}),
	)

	if useDebug {
		_ = logger.Init(logger.WithLevel(logger.DebugLevel))
	}

	_ = servicemgr.NewServiceFactory(microService)

	database, err := db.Connect(config.Config().Database)
	if err != nil {
		logger.Fatalf("Connect to database failed: %s", err)
	}

	camFactory := camera.NewFactory()
	settingsProvider := settings.New()
	cctvService := service.Service{
		Database:         database,
		CameraManager:    manager.New(camFactory, cctv.New()),
		Reactor:          reactor.New(),
		Notifier:         pubsub.NewPublisher(microService),
		ReactFactory:     reactions.NewFactory(pubsub.NewPublisher(microService), settingsProvider),
		SettingsProvider: settingsProvider,
		Timeline:         timeline.New(),
	}

	if err = cctvService.Initialize(); err != nil {
		logger.Fatalf("Initialize service failed: %s", err)
	}

	// регистрируем хендлеры
	if err = rms_cctv.RegisterRmsCctvHandler(microService.Server(), &cctvService); err != nil {
		logger.Fatalf("Register microService failed: %s", err)
	}

	if err = microService.Run(); err != nil {
		logger.Fatalf("Run service failed: %s", err)
	}
}
