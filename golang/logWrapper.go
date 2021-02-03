package golang

//LoggerWrapper - обертка, адаптер для любого типа логера
type LoggerWrapper func(string, ...interface{})

// Printf - реализация интерфейса необходимого для работы retryablehttp
func (log LoggerWrapper) Printf(str string, args ...interface{}) {
	log(str, args...)
}
