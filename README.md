# VK

## Установка

1. Установите Go
2. Склонируйте репозиторий: git clone https://github.com/yourusername/mattermost-voting-bot.git
cd mattermost-voting-bot
3. Установите зависимости: go mod tidy
4. Установите свой Token и необходимые параметры для бота в config.go
5. Запустите Tarantool отдельно (если не используете Docker): tarantool
6. Запустите бота: go run cmd/main.go
7. Либо используйте Docker: docker-compose up --build

## Команды

1. /createpoll Вопрос — создать голосование
2. /vote ID Вариант — проголосовать
3. /results ID — посмотреть результаты
4. /closepoll ID — закрыть голосование
5. /deletepoll ID — удалить голосование
