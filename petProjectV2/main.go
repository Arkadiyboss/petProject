package main

import (
	"fmt"
	"net/http"
	"petProjectV2/config"
	"petProjectV2/internal/transport"
	"petProjectV2/pkg"
)

func main() {

	fmt.Println(123321)

	config := config.LoadConfig()

	connString := fmt.Sprint(config.Db.Adress + ":" + config.Secret.DbSecret + "@db" + config.Db.Port + "/mydb?sslmode=" + config.Db.Sslmode)

	fmt.Println(connString)

	pool, err := pkg.ConnectDB(connString)

	if err != nil {
		fmt.Println("Не получается подключиться к БД", err)
		return
	}

	defer pool.Close()

	go pkg.PingDB(pool, connString)

	transport.StartServer(pool, config)

	fmt.Println("Сервер запущен")
	http.ListenAndServe(config.Server.Port, nil)
}
