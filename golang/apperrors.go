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
	return fmt.Sprintf("Error Code %d, Description: %s", e.Code, e.Description)
}

func (e EgeonError) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"Code\":%d, \"Description\": %q}", e.Code, e.Description)), nil
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
	BadInsertAttempt
	BadUpdateAttempt
	DoNotBeHere
	Permission
	ValidateError
	ServiceWorkError
)

type errType error

func GetErr(errMsg string) error {
	return errors.New(errMsg)
}

// Errors - карта ошибок
var Errors map[string]errType

func init() {
	Errors = make(map[string]errType)
	Errors["cacheInit"] = GetErr("Кеш не ініційовано")
	Errors["badUser"] = GetErr("Користувач має не коректний тип")
	Errors["undefUser"] = GetErr("Невідомий користувач. Ця операція вимагає авторизації користувача")
	Errors["undefDev"] = GetErr("Модем не знайдено")
	Errors["badCMD"] = GetErr("Команда до пристрою не коректна, або не підримується пристроєм")
	Errors["badType"] = GetErr("Вхідні данні мають не вірний тип. Або показчик в нікуди. Зверніться до адміністратора")
	Errors["permission"] = GetErr("У Вас немає прав на виконання даної операції. Або час Вашої сесії закінчився")
	Errors["notFound"] = GetErr("Упс. Я не розумію що Ви намагаєтесь зробити. Цей ресурс не доступний, перевірте правильність ведених даних")
	Errors["notImplement"] = GetErr("Цей функціонал не реалізовано")
	Errors["notAllowed"] = GetErr("Цей метод або дія не допустимі. До ресурсу доступ обмежений, або у Вас нема необхідних прав, або час Вашої сесії закінчився")
	Errors["incorrectInput"] = GetErr("Не коректні вхідні данні")
	Errors["notFindRecord"] = GetErr("Прийшов пустий список записів. Таких записів не знайдено")
	Errors["deleteWithotID"] = GetErr("Can not delete value without ID")
	Errors["undefCmd"] = GetErr("Undefine command")
	Errors["badEmail"] = GetErr("Перевірте правильність ведення электроної пошти")
	Errors["badReqID"] = GetErr("Не відомий request ID")
	Errors["badKey"] = GetErr("Не вірний пароль або ключ авторизації")
}
