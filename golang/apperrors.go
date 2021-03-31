package golang

import (
	"errors"
	"fmt"
)

// EgeonError - implement error and json interfaces for errors in system. TODO optimize it later
type EgeonError struct {
	Code        uint32 `json:"Code"`
	Description string `json:"Description"`
}

func (e EgeonError) Error() string {
	return fmt.Sprintf("{\"Code\":%d, \"Description\": \"%s\"}", e.Code, e.Description)
}

func (e EgeonError) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"Code\":%d, \"Description\": \"%s\"}", e.Code, e.Description)), nil
}

const (
	Inserted uint32 = iota
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
	DoNotBeHere
	Permission
	ValidateError
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
	Errors["badType"] = GetErr("Вхідні данні мають не вірний тип. Або покзчик в нікуди. Зверніться до адміністратора")
	Errors["permission"] = GetErr("У Вас немає прав на виконання даної операції")
	Errors["notFound"] = GetErr("Упс. Я не розумію що Ви намагаєтесь зробити. Цей ресурс не доступний")
	Errors["notImplement"] = GetErr("Цей функціонал не реалізовано")
	Errors["notAllowed"] = GetErr("Цей метод або дія не допустимі. До ресурсу доступ обмежений, або у Вас нема необхідних прав")
	Errors["incorrectInput"] = GetErr("Не коректні вхідні данні")
	Errors["notFindRecord"] = GetErr("Прийшов пустий список записів. Таких записів не знайдено")
	Errors["deleteWithotID"] = GetErr("Can not delete value without ID")
	Errors["undefCmd"] = GetErr("Undefine command")
	Errors["badEmail"] = GetErr("Не корректно задана электронная почта")
	Errors["badReqID"] = GetErr("Не известный request ID")
	Errors["badKey"] = GetErr("Не корректный пароль или ключ идентификации")
}
