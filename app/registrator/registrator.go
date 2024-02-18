package main

import (
	"context"
	rms_cctv "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-cctv"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"

	// Plugins
	_ "github.com/go-micro/plugins/v4/registry/etcd"
)

func main() {
	camera := rms_cctv.Camera{Schedule: "{}"}
	service := micro.NewService(
		micro.Name("rms-cctv.client"),
		micro.Flags(
			&cli.StringFlag{
				Name:        "name",
				Usage:       "Camera name",
				Required:    true,
				Destination: &camera.Name,
			},
			&cli.StringFlag{
				Name:        "url",
				Usage:       "Camera URL",
				Required:    true,
				Destination: &camera.Url,
			},
			&cli.StringFlag{
				Name:        "user",
				Usage:       "Username",
				Required:    false,
				Destination: &camera.User,
			},
			&cli.StringFlag{
				Name:        "password",
				Usage:       "Password",
				Required:    false,
				Destination: &camera.Password,
			},
		),
	)
	service.Init()

	client := rms_cctv.NewRmsCctvService("rms-cctv", service.Client())
	_, err := client.AddCamera(context.Background(), &camera)
	if err != nil {
		panic(err)
	}
}
