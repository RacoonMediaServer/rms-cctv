package config

import "github.com/RacoonMediaServer/rms-packages/pkg/configuration"

// Configuration represents entire service configuration
type Configuration struct {
	Database configuration.Database
	Debug    configuration.Debug
	Cctv     Cctv
}

type Backend struct {
	Type  string
	Host  string
	Path  string
	Token string
}

type Cctv struct {
	Backend Backend
}

var config Configuration

// Load open and parses configuration file
func Load(configFilePath string) error {
	return configuration.Load(configFilePath, &config)
}

// Config returns loaded configuration
func Config() Configuration {
	return config
}
