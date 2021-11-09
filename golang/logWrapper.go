package golang

import (
	"fmt"
	"os"
)

// //LoggerWrapper - обертка, адаптер для любого типа логера
// type LoggerWrapper func(string, ...interface{})

// // Printf - реализация интерфейса необходимого для работы retryablehttp
// func (log LoggerWrapper) Printf(str string, args ...interface{}) {
// 	log(str, args...)
// }

type LoggerWriterWraper func(string)

func (log LoggerWriterWraper) Write(buf []byte) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			os.Stderr.WriteString(fmt.Sprintf("%v", err))
		}
	}()
	log(string(buf))
	return len(buf), nil
}

func (log LoggerWriterWraper) WriteString(buf string) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			os.Stderr.WriteString(fmt.Sprintf("%v", err))
		}
	}()
	log(buf)
	return len(buf), nil
}
