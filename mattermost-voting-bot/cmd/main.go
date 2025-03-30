package main

import (
	"log"
	"mattermost-voting-bot/config"
	"mattermost-voting-bot/internal/bot"
	"mattermost-voting-bot/internal/storage"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к Tarantool
	tarantoolConn, err := storage.ConnectTarantool(cfg.TarantoolAddress)
	if err != nil {
		log.Fatalf("Ошибка подключения к Tarantool: %v", err)
	}
	defer tarantoolConn.Close()

	// Запускаем бота
	go bot.Run(cfg, tarantoolConn)

	// Запускаем WebSocket
	bot.StartWebSocket(cfg, tarantoolConn)
}
