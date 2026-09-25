package transport

type Errors struct {
	code int
	text string
}

var httpError = make(map[int]Errors)

func MakeErrors() {

	httpError[400] = Errors{code:400, text: "неправильный, некорректный запрос"}
	httpError[4001] = Errors{code:400, text: "Неверный метод запроса"}
	httpError[4002] = Errors{code:400, text: "Login не заполнен"}
	httpError[4003] = Errors{code:400, text: "Password не заполнен"}
	httpError[4004] = Errors{code:400, text: "Пользователь уже существует"}
	httpError[4005] = Errors{code:403, text: "Неверный пароль пользователя"}
	httpError[4006] = Errors{code:403, text: "Авторизационный токен пустой"}
	httpError[4007] = Errors{code:403, text: "Токен истек, залогиньтесь заново"}
	httpError[404] = Errors{code:404, text: "Пользователь не найден"}
	httpError[500] = Errors{code:500, text: "Внутренняя ошибка сервера"}
	httpError[5001] = Errors{code:500, text: "Ошибка парсинга логина и пароля"}
	httpError[5002] = Errors{code:500, text: "Ошибка при шифровании пароля"}
	httpError[5003] = Errors{code:500, text: "Ошибка записи в БД"}
	httpError[5004] = Errors{code:500, text: "Ошибка при создании токена"}
	httpError[5005] = Errors{code:500, text: "Ошибка при поиске токена"}
	httpError[5006] = Errors{code:500, text: "Ошибка при попытке обновления токена"}
	httpError[5007] = Errors{code:500, text: "Ошибка парсинга входящего токена"}

}

func FindError(errorCode int) (int, string){
	httpError := httpError[errorCode]
	return httpError.code, httpError.text
}