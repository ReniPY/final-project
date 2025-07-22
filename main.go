package main

import (
	"fmt"

	"github.com/ReniPY/final-project/pkg/db"
	"github.com/ReniPY/final-project/pkg/server"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		fmt.Println("Ошибка инициализации базы данных:", err)
		return
	}

	server.Run()
}
