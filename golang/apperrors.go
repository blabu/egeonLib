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
	Errors["cacheInit"] = GetErr("Кеш не инициализирован")
	Errors["badUser"] = GetErr("Не корректный тип пользователя")
	Errors["undefUser"] = GetErr("Пользователь не известен. Это действие требует авторизации пользователя")
	Errors["undefDev"] = GetErr("Такой модем не найден")
	Errors["badCMD"] = GetErr("Не корректная команда, или команда не поддерживается модемом")
	Errors["badType"] = GetErr("Входные данные имеют не корректный тип")
	Errors["permission"] = GetErr("У Вас нет прав на осуществление данной операции")
	Errors["notFound"] = GetErr("Упс. Я не понимаю чего Вы хотите")
	Errors["notImplement"] = GetErr("Method not implemented")
	Errors["notAllowed"] = GetErr("Этот метод или дейсвие не допустимо. К запрашиваемому ресурсу доступ ограничен")
	Errors["incorrectInput"] = GetErr("Не корректные входные значения")
	Errors["notFindRecord"] = GetErr("Ни одной записи не найдено")
	Errors["deleteWithotID"] = GetErr("Can not delete value without ID")
	Errors["undefCmd"] = GetErr("Undefine command")
	Errors["badEmail"] = GetErr("Не корректно задана электронная почта")
	Errors["badReqID"] = GetErr("Не известный request ID")
	Errors["badKey"] = GetErr("Не корректный пароль или ключ идентификации")
}
