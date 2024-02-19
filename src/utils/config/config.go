package config

import "go.uber.org/zap"

type ServerConfig struct {
	ENV         string
	PORT        string
	MTRG_CLIENT string
	LOGGER      *zap.Logger
}

type Router struct {
	Name         string
	Path         string
	ServerConfig *ServerConfig
}
