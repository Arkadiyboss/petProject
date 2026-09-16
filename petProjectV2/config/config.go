package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db     db	`json:"db"`
	Server server	`json:"server"`
	Jwt    jwt	`json:"jwt"`
	Password password `json:"password"`
	Secret secret	`json:"secret"`
}

type db struct {
	Adress  string `json:"adress"`
	Port    string `json:"port"`
	Sslmode string `json:"sslmode"`
}

type server struct {
	Port string	`json:"port"`
}

type jwt struct {
	LenghtPass int	`json:"lenghtPass"`
}

type password struct {
	Key string `json:"key"`
}

type secret struct {
	JwtSecret string `json:"jwtSecret"`
	DbSecret string `json:"dbSecret"`
}

type readSecretPath struct {
	JwtSecretName string	`json:"jwtName"`
	DbSecretName  string `json:"dbName"`
}



func LoadConfig() *Config{

	newConfig := Config{}

	byteConfig, err := os.ReadFile("config.json")

	if err != nil {
		fmt.Println("Ошибка чтения файла: ")
		panic(err)
	}

	json.Unmarshal(byteConfig, &newConfig)

	newConfig.Secret.JwtSecret, newConfig.Secret.DbSecret, err = ReadSecret()

	if err != nil {
		fmt.Println("Ошибка получения секрета: ")
		panic(err)
	}

	return &newConfig
}

func ReadSecret() (string, string, error){

	newSecret := readSecretPath{}

	byteSecret, err :=os.ReadFile("secret.json")

	if err != nil {
		fmt.Println("Ошибка чтения секретов: ", err)
		return "", "", err
	}

	err = json.Unmarshal(byteSecret, &newSecret)

	if err != nil {
		fmt.Println("Ошибка парсинга секретов: ", err)
		return "","", err
	}

	JwtSecret, DbSecret := ReadEnv(&newSecret)

	
	return JwtSecret, DbSecret, nil

}

func ReadEnv(newSecret *readSecretPath) (string, string){



	JwtSecret := os.Getenv(newSecret.JwtSecretName)

	if JwtSecret == "" {
		panic("Не задан jwt секрет")
	}

	DbSecret := os.Getenv(newSecret.DbSecretName)

	if DbSecret == "" {
		panic("Не задан пароль БД")
	}


	return JwtSecret, DbSecret

}
