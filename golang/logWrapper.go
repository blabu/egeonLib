package golang

//LoggerWrapper - обертка, адаптер для любого типа логера
type LoggerWrapper func(string, ...interface{})

// Printf - реализация интерфейса необходимого для работы retryablehttp
func (log LoggerWrapper) Printf(str string, args ...interface{}) {
	log(str, args...)
}

type LoggerWriterWraper func(string)

func (log LoggerWriterWraper) Write(buf []byte) (int, error) {
	log(string(buf))
	return len(buf), nil
}

func (log LoggerWriterWraper) WriteString(buf string) (int, error) {
	log(buf)
	return len(buf), nil
}
