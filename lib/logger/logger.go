package logger

import (
	"fmt"
	"time"
)

func Log(m string) {
	fmt.Println(time.Now().Format(time.RFC3339) + " | LOG | " + m)
}

func Error(err error) {
	fmt.Println(time.Now().Format(time.RFC3339) + " | FATAL | " + err.Error())
}
