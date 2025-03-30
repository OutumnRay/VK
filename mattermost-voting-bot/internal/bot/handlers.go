package bot

import (
	"fmt"
	"log"

	"mattermost-voting-bot/config"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/tarantool/go-tarantool"
)

func processCommand(post *model.Post, client *model.Client4, cfg config.Config, tarantoolConn *tarantool.Connection) {
	message := strings.ToLower(post.Message)

	switch {
	case strings.HasPrefix(message, "/createpoll"):
		createPoll(post, client, tarantoolConn)
	case strings.HasPrefix(message, "/vote"):
		votePoll(post, client, tarantoolConn)
	case strings.HasPrefix(message, "/results"):
		showResults(post, client, tarantoolConn)
	case strings.HasPrefix(message, "/closepoll"):
		closePoll(post, client, tarantoolConn)
	case strings.HasPrefix(message, "/deletepoll"):
		deletePoll(post, client, tarantoolConn)
	}
}

// createPoll создает голосование
func createPoll(post *model.Post, client *model.Client4, tarantoolConn *tarantool.Connection) {
	args := strings.Split(post.Message, "\n")
	if len(args) < 3 {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: недостаточно аргументов. Используйте:\n/create_poll Вопрос\nВариант 1\nВариант 2\n...",
		})
		return
	}

	question := args[1]
	options := args[2:]

	pollID := fmt.Sprintf("%s_%s", post.Id, post.UserId)
	poll := map[string]interface{}{
		"id":       pollID,
		"question": question,
		"options":  options,
		"votes":    map[string]int{},
		"creator":  post.UserId,
	}

	_, err := tarantoolConn.Insert("polls", []interface{}{pollID, poll})
	if err != nil {
		log.Println("Ошибка сохранения голосования в Tarantool:", err)
		return
	}

	client.CreatePost(&model.Post{
		ChannelId: post.ChannelId,
		Message:   fmt.Sprintf("Голосование создано! ID: %s\n%s\n%s", pollID, question, strings.Join(options, "\n")),
	})
}

// votePoll позволяет пользователям голосовать
func votePoll(post *model.Post, client *model.Client4, tarantoolConn *tarantool.Connection) {
	args := strings.Fields(post.Message)
	if len(args) < 3 {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: недостаточно аргументов. Используйте:\n/vote ID_Голосования Вариант",
		})
		return
	}

	pollID, option := args[1], strings.Join(args[2:], " ")

	var poll map[string]interface{}
	err := tarantoolConn.GetTyped("polls", "primary", []interface{}{pollID}, &poll)
	if err != nil {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: голосование не найдено!",
		})
		return
	}

	votes := poll["votes"].(map[string]int)
	votes[option]++

	_, err = tarantoolConn.Replace("polls", []interface{}{pollID, poll})
	if err != nil {
		log.Println("Ошибка обновления голосования:", err)
	}

	client.CreatePost(&model.Post{
		ChannelId: post.ChannelId,
		Message:   "Ваш голос учтен!",
	})
}

// showResults показывает результаты голосования
func showResults(post *model.Post, client *model.Client4, tarantoolConn *tarantool.Connection) {
	args := strings.Fields(post.Message)
	if len(args) < 2 {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: укажите ID голосования.\n/show_results ID_Голосования",
		})
		return
	}

	pollID := args[1]

	var poll map[string]interface{}
	err := tarantoolConn.GetTyped("polls", "primary", []interface{}{pollID}, &poll)
	if err != nil {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: голосование не найдено!",
		})
		return
	}

	results := "Результаты голосования:\n"
	for option, count := range poll["votes"].(map[string]int) {
		results += fmt.Sprintf("%s: %d голосов\n", option, count)
	}

	client.CreatePost(&model.Post{
		ChannelId: post.ChannelId,
		Message:   results,
	})
}

// closePoll завершает голосование
func closePoll(post *model.Post, client *model.Client4, tarantoolConn *tarantool.Connection) {
	args := strings.Fields(post.Message)
	if len(args) < 2 {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: укажите ID голосования.\n/close_poll ID_Голосования",
		})
		return
	}

	pollID := args[1]

	var poll map[string]interface{}
	err := tarantoolConn.GetTyped("polls", "primary", []interface{}{pollID}, &poll)
	if err != nil {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: голосование не найдено!",
		})
		return
	}

	if poll["creator"].(string) != post.UserId {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: завершать голосование может только создатель.",
		})
		return
	}

	_, err = tarantoolConn.Delete("polls", "primary", []interface{}{pollID})
	if err != nil {
		log.Println("Ошибка удаления голосования:", err)
	}

	client.CreatePost(&model.Post{
		ChannelId: post.ChannelId,
		Message:   "Голосование завершено!",
	})
}

// deletePoll удаляет голосование
func deletePoll(post *model.Post, client *model.Client4, tarantoolConn *tarantool.Connection) {
	args := strings.Fields(post.Message)
	if len(args) < 2 {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: укажите ID голосования.\n/delete_poll ID_Голосования",
		})
		return
	}

	pollID := args[1]

	var poll map[string]interface{}
	err := tarantoolConn.GetTyped("polls", "primary", []interface{}{pollID}, &poll)
	if err != nil {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: голосование не найдено!",
		})
		return
	}

	if poll["creator"].(string) != post.UserId {
		client.CreatePost(&model.Post{
			ChannelId: post.ChannelId,
			Message:   "Ошибка: удалять голосование может только создатель.",
		})
		return
	}

	_, err = tarantoolConn.Delete("polls", "primary", []interface{}{pollID})
	if err != nil {
		log.Println("Ошибка удаления голосования:", err)
	}

	client.CreatePost(&model.Post{
		ChannelId: post.ChannelId,
		Message:   "Голосование удалено!",
	})
}
