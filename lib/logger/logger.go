package logger

import (
	"fmt"
	"time"
)

func Log(m string) {
	fmt.Println(time.Now().Format(time.RFC3339) + " | " + m)
}

func Fatal(m string) {
	fmt.Println(time.Now().Format(time.RFC3339) + " | " + m)
}
