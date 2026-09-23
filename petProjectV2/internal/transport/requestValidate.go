package transport

import (
	"encoding/json"
	"net/http"
	"petProjectV2/config"
	auth "petProjectV2/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func LoginPasswordValidate(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, config *config.Config, endPoint string) {

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

	respondInterface := auth.NewRequest()

	switch endPoint {
	case "/register":

		respond, errCode := respondInterface.UserRegister(login, password, config, pool)

		if errCode != 0 {
			httpErrCode, errText := FindError(errCode)
			http.Error(w, errText, httpErrCode)
			return
		}

		w.Header().Set("Authorization", respond.AccessToken)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(respond.Message)

	case "/login":

		respond, errCode := respondInterface.UserLogin(login, password, config, pool)

		if errCode != 0 {
			httpErrCode, errText := FindError(errCode)
			http.Error(w, errText, httpErrCode)
			return
		}

		w.Header().Set("Authorization", respond.AccessToken)
		w.Header().Set("Allow", respond.RefreshToken)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(respond.Message)
	}

}

func TokenValidate(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, config *config.Config) {

	respondInterface := auth.NewRequest()

	if r.Method != http.MethodGet {
		httpErrCode, errText := FindError(4001)
		http.Error(w, errText, httpErrCode)
		return
	}

	token := r.Header.Get("Authorization")

	if token == "" {
		httpErrCode, errText := FindError(4006)
		http.Error(w, errText, httpErrCode)
		return
	}

	respond, errCode := respondInterface.Verify(token, config, pool)

	if errCode != 0 {
		httpErrCode, errText := FindError(errCode)
		http.Error(w, errText, httpErrCode)
		return
	}

	w.Header().Set("Authorization", respond.AccessToken)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respond.Message)

}
