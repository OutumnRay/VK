package bot

import (
	"log"

	"mattermost-voting-bot/config"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/tarantool/go-tarantool"
)

// Run запускает бота
func Run(cfg config.Config, tarantoolConn *tarantool.Connection) {
	client := model.NewAPIv4Client(cfg.MattermostServer)
	client.SetToken(cfg.MattermostToken)

	// Подключаем WebSocket
	wsClient, err := model.NewWebSocketClient4(cfg.MattermostServer, client.AuthToken)
	if err != nil {
		log.Fatalf("Ошибка подключения к WebSocket: %v", err)
	}
	wsClient.Listen()
	defer wsClient.Close()

	// Обрабатываем события
	for event := range wsClient.EventChannel {
		go handleEvent(event, client, cfg, tarantoolConn)
	}
}
