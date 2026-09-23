package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"petProjectV2/config"
	"strings"
	"time"

	pass "petProjectV2/internal/service/password"
	"petProjectV2/pkg"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// 4001 - Неверный метод запроса
// 4002 - Login не заполнен
// 4003 - Password не заполнен
// 4004 - Пользователь уже существует
// 4005 - Неверный пароль пользователя
// 4006 - Авторизационный токен пустой
// 4007 - Токен истек, залогиньтесь заново
// 5001 - Ошибка парсинга логина и пароля
// 5002 - Ошибка при шифровании пароля
// 5003 - Ошибка записи в БД
// 5004 - Ошибка при создании токена
// 5005 - Ошибка при попытке поиска токена
// 5006 - Ошибка при попытке обновления токена
// 5007 - Ошибка парсинга входящего токена

type UserService interface {
    UserRegister(login, password string, config *config.Config, pool *pgxpool.Pool) (*RequestResult, int)
    UserLogin(login, password string, config *config.Config, pool *pgxpool.Pool) (*RequestResult, int)
    Verify(token string, config *config.Config, pool *pgxpool.Pool) (*RequestResult, int)
}

type RequestResult struct {
    Message string
    AccessToken   string
	RefreshToken   string
}

type Token struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}


func (req *RequestResult) Verify(token string, config *config.Config, pool *pgxpool.Pool) (*RequestResult, int) {

	valid, errCode := DecodeAccessToken(token)

	if errCode != 0 {
		return nil, errCode
	}

	if !valid {
		return nil, 4007
	}

	var dbLogin string
	var dbPass string

	dbLogin, dbPass, errCode = pkg.VerifyToken(token, pool)

	if errCode != 0 {
		return nil, errCode
	}

	updatedTimeToken, errCode := AccessToken(dbLogin, dbPass, config)

	if errCode != 0 {
		return nil, errCode
	}

	return &RequestResult{
	AccessToken: updatedTimeToken,}, errCode

}

func AccessToken(login, password string, config *config.Config) (string, int) {
	claims := Token{
		Login:     login,
		Password:  password,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Secret.JwtSecret))

	if err != nil {
		return "", 5004
	}

	return tokenString, 0
}

func DecodeAccessToken(pass string) (bool, int) {

	part := strings.Split(pass, ".")

	body, err := base64.RawURLEncoding.DecodeString(part[1])

	if err != nil {
		return false, 5007
	}

	var info map[string]interface{}

	err = json.Unmarshal(body, &info)

	if err != nil {
		return false, 5007
	}

	expired := info["exp"].(float64)

	if time.Now().Unix() > int64(expired) {
		return false, 5007
	}

	return true, 0

}

func RefreshToken(config *config.Config) string {
	avaliableSymbols := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	var pass string
	for i := 0; i <= config.Jwt.LenghtPass-1; i++ {
		symbol := rand.Intn(cap(avaliableSymbols))
		pass = pass + avaliableSymbols[symbol]
	}
	return pass
}

func (r *RequestResult)UserRegister(login string, password string, config *config.Config, pool *pgxpool.Pool) (*RequestResult, int) {


	query := `Select id From users Where "login" = $1`

	id, err := pkg.FindOrdinaryInt(query, login, pool)

	if err == nil {
		return nil, 4004
	}

	if id != 0 {
		return nil, 4004
	}

	envPass, err := pass.EncryptPass(password, config)

	if err != nil {
		return nil, 5002
	}

	id, err = pkg.RegisterUser(login, envPass, pool)

	if err != nil {
		return nil, 5003
	}

	successText := fmt.Sprint("Пользователь успешно создан, его айди: ", id)

	accessToken, errCode := AccessToken(login, password, config)

	if errCode != 0 {
		return nil, 5003
	}

	return &RequestResult{Message: successText,
		AccessToken: accessToken,}, 0
}

func (r *RequestResult)UserLogin(login string, password string, config *config.Config, pool *pgxpool.Pool) (*RequestResult, int) {


	query := `Select id From users Where "login" = $1`

	id, err := pkg.FindOrdinaryInt(query, login, pool)

	if err != nil {
		return nil, 404
	}

	query = `Select "passwordHash" From users Where "login" = $1`

	dbPass, err := pkg.FindOrdinaryString(query, login, pool)

	if err != nil {
		return nil, 500
	}

	if id == 0 {
		return nil, 404
	}

	pass, err := pass.DecryptPass(dbPass, config)

	if err != nil {
		return nil, 5002
	}

	if pass != password {
		return nil, 4005
	}

	refresh := RefreshToken(config)

	access, errCode := AccessToken(login, dbPass, config)

	if errCode != 0 {
		return nil, 5004
	}

	success := pkg.UpdateSession(access, id, pool)

	if !success {
		return nil, 5003
	}

	successText := "Логин выполнен успешно"

	return &RequestResult{Message: successText,
	AccessToken: access,
	RefreshToken: refresh,}, 0

}

func NewRequest() (*RequestResult){
	return &RequestResult{}
}