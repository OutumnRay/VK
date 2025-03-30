package bot

import (
	"encoding/json"
	"log"

	"mattermost-voting-bot/config"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/tarantool/go-tarantool"
)

// StartWebSocket подключается к WebSocket и начинает слушать события
func StartWebSocket(cfg config.Config, tarantoolConn *tarantool.Connection) {
	client := model.NewAPIv4Client(cfg.MattermostServer)
	client.SetToken(cfg.MattermostToken)

	wsClient, err := model.NewWebSocketClient4(cfg.MattermostServer, client.AuthToken)
	if err != nil {
		log.Fatalf("Ошибка подключения к WebSocket: %v", err)
	}
	defer wsClient.Close()

	wsClient.Listen()

	log.Println("WebSocket подключен, ожидаем события...")

	for event := range wsClient.EventChannel {
		go handleEvent(event, client, cfg, tarantoolConn)
	}
}

func handleEvent(event *model.WebSocketEvent, client *model.Client4, cfg config.Config, tarantoolConn *tarantool.Connection) {
	if event.GetBroadcast().ChannelId != cfg.MattermostChannel {
		return
	}

	if event.EventType() != model.WebsocketEventPosted {
		return
	}

	var post model.Post
	if err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post); err != nil {
		log.Println("Ошибка обработки сообщения:", err)
		return
	}

	// Игнорируем сообщения от самого бота
	if post.UserId == cfg.MattermostToken {
		return
	}

	processCommand(&post, client, cfg, tarantoolConn)
}
