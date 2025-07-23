package main

import (
	"fmt"

	"github.com/ReniPY/final-project/pkg/db"
	"github.com/ReniPY/final-project/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		fmt.Println("Ошибка инициализации базы данных:", err)
		return
	}
	defer db.Close()

	server.Run()
}
