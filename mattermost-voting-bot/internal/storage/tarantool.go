package storage

import (
	"log"

	"github.com/tarantool/go-tarantool"
)

// ConnectTarantool подключается к Tarantool
func ConnectTarantool(address string) (*tarantool.Connection, error) {
	conn, err := tarantool.Connect(address, tarantool.Opts{User: "guest"})
	if err != nil {
		return nil, err
	}
	log.Println("Подключено к Tarantool")
	return conn, nil
}

// Функции для сохранения, обновления и удаления голосований...
