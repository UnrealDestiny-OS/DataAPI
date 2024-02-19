package config

import "go.uber.org/zap"

type ServerConfig struct {
	ENV         string
	PORT        string
	MTRG_CLIENT string
	LOGGER      *zap.Logger
}
