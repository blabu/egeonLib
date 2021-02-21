package golang

import "errors"

const (
	Inserted uint16 = iota
	Shared
	HistoryChanged
	Accepted
	Deleted
	Updated
	NotAuthError
	PageNotFind
	NotFindItemError
	IncorrectRequestParam
	DatabaseError
	ReqToBig
	ReportError
	InternalError
	Wait
	MethodNotImplemented
	BadDeleteAttempt
)

type errType error

func GetErr(errMsg string) error {
	return errors.New(errMsg)
}

//Errors - карта ошибок
var Errors map[string]errType

func init() {
	Errors = make(map[string]errType)
	Errors["cacheInit"] = GetErr("Кеш не ініційовано")
	Errors["badUser"] = GetErr("Користувач має не коректний тип")
	Errors["undefUser"] = GetErr("Невідомий користува. Ця операція вимагає авторизації користувача")
	Errors["undefDev"] = GetErr("Модем не знайдено")
	Errors["badCMD"] = GetErr("Команда до пристрою не коректна, або не підримується пристроєм")
	Errors["badType"] = GetErr("Вхідні данні мають не вірний тип")
	Errors["permission"] = GetErr("У Вас немає прав на виконання даної операції")
	Errors["notFound"] = GetErr("Упс. Я не розумію що Ви намагаєтесь зробити. Цей ресурс не доступний")
	Errors["notImplement"] = GetErr("Цей функціонал не реалізовано")
	Errors["notAllowed"] = GetErr("Цей метод або дія не лопустимі. До ресурсу доступ обмежений")
	Errors["incorrectInput"] = GetErr("Не коректні вхідні данні")
	Errors["notFindRecord"] = GetErr("Прийшов пустий список записів. Таких записів не знайдено")
	Errors["deleteWithotID"] = GetErr("Can not delete value without ID")
	Errors["undefCmd"] = GetErr("Undefine command")
	Errors["badEmail"] = GetErr("Не корректно задана электронная почта")
	Errors["badReqID"] = GetErr("Не известный request ID")
	Errors["badKey"] = GetErr("Не корректный пароль или ключ идентификации")
}
