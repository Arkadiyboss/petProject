package pkg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FindOrdinaryInt(query string, value any, pool *pgxpool.Pool) (int, error){

	var res int

	err := pool.QueryRow(context.Background(), query, value).Scan(&res)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, err
		}
		return 0, err
	}
	return res, nil
}

func FindOrdinaryString(query string, value any, pool *pgxpool.Pool) (string, error){

	var res string

	err := pool.QueryRow(context.Background(), query, value).Scan(&res)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", err
		}
		return "0", err
	}
	return res, nil
}

func FindMany() {

}

func RegisterUser(login string, envPass string, pool *pgxpool.Pool) (int, error){

	query := `INSERT INTO users (login, "passwordHash") VALUES ($1, $2) RETURNING id`

		var id int

		err := pool.QueryRow(
			context.Background(),
			query,
			login,
			envPass,
		).Scan(&id)

		if err != nil {
			return 0, err
		}

		return id, nil
}

func UpdateSession(access string, id int, pool *pgxpool.Pool) bool {

	query := `INSERT INTO session (access_token, user_id) VALUES ($1, $2)`

		_, err := pool.Exec(
			context.Background(),
			query,
			access,
			id,
		)

		if err != nil {
			return false
		}
		return true
}

func VerifyToken(token string, pool *pgxpool.Pool) (string, string, int){

	var dbLogin string
	var dbPass string

	query := `SELECT login, "passwordHash" from users u JOIN session s ON u.id = s.user_id  WHERE access_token = $1`

	err := pool.QueryRow(context.Background(), query, token).Scan(&dbLogin, &dbPass)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", "", 404
		}
		return "", "", 5005
	}

	return dbLogin, dbPass, 0
}
