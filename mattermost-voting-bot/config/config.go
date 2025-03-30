package config

import (
	"os"
)

// Config хранит параметры бота
type Config struct {
	MattermostServer  string
	MattermostToken   string
	MattermostTeam    string
	MattermostChannel string
	TarantoolAddress  string
}

// Load загружает конфигурацию из переменных окружения
func Load() Config {
	return Config{
		MattermostServer:  os.Getenv("MATTERMOST_SERVER"),
		MattermostToken:   os.Getenv("MATTERMOST_TOKEN"),
		MattermostTeam:    os.Getenv("MATTERMOST_TEAM"),
		MattermostChannel: os.Getenv("MATTERMOST_CHANNEL"),
		TarantoolAddress:  os.Getenv("TARANTOOL_ADDRESS"),
	}
}
