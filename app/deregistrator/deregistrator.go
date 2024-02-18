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
	id := uint(0)
	service := micro.NewService(
		micro.Name("rms-cctv.client"),
		micro.Flags(
			&cli.UintFlag{
				Name:        "id",
				Usage:       "Camera ID",
				Required:    true,
				Destination: &id,
			},
		),
	)
	service.Init()

	client := rms_cctv.NewRmsCctvService("rms-cctv", service.Client())
	_, err := client.DeleteCamera(context.Background(), &rms_cctv.DeleteCameraRequest{CameraId: uint32(id)})
	if err != nil {
		panic(err)
	}
}
