package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LogInfo struct {
	Error error     `json:"error"`
	Time  time.Time `json:"time"`
}

func ConnectDB(connStr string) (*pgxpool.Pool, error) {

	Pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		fmt.Println("Не получается подключиться к БД", err)
		return Pool, err
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		fmt.Println("Не получается подключиться к БД", err)
		return Pool, err

	}

	return Pool, err
}

func PingDB(Pool *pgxpool.Pool) {

	timer := time.NewTicker(1 * time.Minute)
	defer timer.Stop()

	for range timer.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := Pool.Ping(ctx)
		cancel()

		if err != nil {
			_, err := json.Marshal(Pool)
			if err != nil {
				fmt.Println("Не смог перевести байты в лог")
			}
			ConnectDB("postgres://arkadiy:123@localhost:5432/mydb?sslmode=disable")
		}
	}
}
