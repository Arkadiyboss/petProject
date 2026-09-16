package main

import (
	"fmt"
	"net/http"
	"petProjectV2/config"
	"petProjectV2/pkg"
	"petProjectV2/internal/service"
)

func main() {
	
	config := config.LoadConfig()

	connString := fmt.Sprint(config.Db.Adress+":"+config.Secret.DbSecret+"@localhost"+config.Db.Port+"/mydb?sslmode="+config.Db.Sslmode)

	fmt.Println(connString)

	pool, err := database.ConnectDB(connString)

	if err != nil {
		fmt.Println("Не получается подключиться к БД", err)
		return
	}

	defer pool.Close()

	go database.PingDB(pool, connString)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginPassword(w, r, pool, config, 1)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginPassword(w, r, pool, config, 2)
	})
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		auth.Verify(w, r, pool, config)
	})

	fmt.Println("Сервер запущен")
	http.ListenAndServe(config.Server.Port, nil)
}
