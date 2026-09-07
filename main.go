package main

import (
	"fmt"
	"net/http"
	"petProject/auth"
	"petProject/database"
)

func main() {
	pool, err := database.ConnectDB("postgres://postgres:123@db:5432/mydb?sslmode=disable")

	if err != nil {
		fmt.Println("Не получается подключиться к БД", err)
		return
	}

	defer pool.Close()

	go database.PingDB(pool)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginPassword(w, r, pool, 1)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginPassword(w, r, pool, 2)
	})
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		auth.Verify(w, r, pool)
	})

	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8080", nil)
}
