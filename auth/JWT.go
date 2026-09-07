package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	pass "petProject/password"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var jwtSecret = []byte("arkadiy-boss")

var lenghtPass int = 48

type Token struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

func LoginPassword(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, number int) {

	if r.Method != http.MethodGet {
		http.Error(w, "Неверный метод запроса", 400)
		return
	}

	login, password, ok := r.BasicAuth()

	if !ok {
		http.Error(w, "Ошибка парсинга логина и пароля", 500)
		return
	}

	if login == "" {
		http.Error(w, "Login не заполнен", 400)
		return
	}
	if password == "" {
		http.Error(w, "Password не заполнен", 400)
		return
	}

	id, dbPass := CheckDB(w, login, pool, number)

	switch number {
	case 1:
		if id != 0 {
			errorText := fmt.Sprint("Пользователь уже существует с айди: ", id)
			http.Error(w, errorText, 400)
			return
		}

		envPass, err := pass.EncryptPass(password)

		if err != nil {
			http.Error(w, "Ошибка при шифровании пароля", 500)
			return
		}

		query := `INSERT INTO users (login, "passwordHash") VALUES ($1, $2) RETURNING id`

		var id int

		err = pool.QueryRow(
			context.Background(),
			query,
			login,
			envPass,
		).Scan(&id)

		if err != nil {
			http.Error(w, "Ошибка записи в БД", 500)
			return
		}

		successText := fmt.Sprint("Пользователь успешно создан, его айди: ", id)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(successText)

	case 2:
		if id == 0 {
			http.Error(w, "Пользователь не найден", 404)
			return
		}

		decryptPass, err := pass.DecryptPass(dbPass)

		if err != nil {
			http.Error(w, "Ошибка расшифровки пароля", 400)
			return
		}

		if password != decryptPass {
			http.Error(w, "Неверный пароль", 400)
			return
		}

		refresh := RefreshToken()

		access, err := AccessToken(login, dbPass)

		if err != nil {
			http.Error(w, "Ошибка при создании токена", 500)
			return
		}

		query := `INSERT INTO session (access_token, user_id) VALUES ($1, $2)`

		_, err = pool.Exec(
			context.Background(),
			query,
			access,
			id,
		)

		if err != nil {
			http.Error(w, "Ошибка записи токена в БД", 500)
			return
		}

		w.Header().Set("Authorization", access)
		w.Header().Set("Allow", refresh)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Токен успешно создан")
	}

}

func Verify(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {

	if r.Method != http.MethodGet {
		http.Error(w, "Неверный метод запроса", 400)
		return
	}

	token := r.Header.Get("Authorization")

	if token == "" {
		http.Error(w, "Авторизационный токен пустой", 400)
		return
	}

	valid, err := DecodeAccessToken(token)

	if err != nil {
		http.Error(w, "Ошибка декодирования токена", 500)
		return
	}

	if !valid {
		http.Error(w, "Токен истек, залогинтесь заново", 403)
		return
	}

	var dbLogin string
	var dbPass string

	query := `SELECT login, "passwordHash" from users u JOIN session s ON u.id = s.user_id  WHERE access_token = $1`


	err = pool.QueryRow(context.Background(), query, token).Scan(&dbLogin, &dbPass)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Входящий токен не найден в БД", 404)
		}
		http.Error(w, "Ошибка при попытке поиска токена", 500)
	}
	decryptPass, err := pass.DecryptPass(dbPass)

	if err != nil {
		http.Error(w, "Ошибка расширофки пароля", 500)
		return
	}

	updatedTimeToken, err := AccessToken(dbLogin, decryptPass)

	if err != nil {
		http.Error(w, "Ошибка обновления токена", 403)
		return
	}

	w.Header().Set("Authorization", updatedTimeToken)
	w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Токен успешно создан")

}

func CheckDB(w http.ResponseWriter, login string, pool *pgxpool.Pool, number int) (int, string) {

	var id int

	query := `Select id From users Where "login" = $1`

	err := pool.QueryRow(context.Background(), query, login).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, ""
		}
		return 0, ""
	}

	var dbPass string
	query = `Select "passwordHash" From users Where "login" = $1`

	err = pool.QueryRow(context.Background(), query, login).Scan(&dbPass)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, ""
		}
		return 0, ""
	}

	return id, dbPass
}

func AccessToken(login, password string) (string, error) {
	claims := Token{
		Login:    login,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeAccessToken(pass string) (bool, error) {

	part := strings.Split(pass, ".")

	body, err := base64.RawURLEncoding.DecodeString(part[1])

	if err != nil {
		return false, err
	}

	var info map[string]interface{}

	err = json.Unmarshal(body, &info)

	if err != nil {
		return false, err
	}

	expired := info["exp"].(float64)

	if time.Now().Unix() > int64(expired) {
		return false, nil
	}

	return true, nil

}

func RefreshToken() string {
	avaliableSymbols := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	var pass string
	for i := 0; i <= lenghtPass-1; i++ {
		symbol := rand.Intn(cap(avaliableSymbols))
		pass = pass + avaliableSymbols[symbol]
	}
	return pass
}
