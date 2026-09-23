package transport

import (
	Http "net/http"
	Config "petProjectV2/config"

	Pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

func StartServer(pool *Pgxpool.Pool, config *Config.Config) {

	MakeErrors()

	Http.HandleFunc("/register", func(w Http.ResponseWriter, r *Http.Request) {
		LoginPasswordValidate(w, r, pool, config, "/register")
	})
	Http.HandleFunc("/login", func(w Http.ResponseWriter, r *Http.Request) {
		LoginPasswordValidate(w, r, pool, config, "/login")
	})
	Http.HandleFunc("/verify", func(w Http.ResponseWriter, r *Http.Request) {
		TokenValidate(w, r, pool, config)
	})

}
