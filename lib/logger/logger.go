package logger

import (
	"fmt"
	"time"
)

func Log(m string) {
	fmt.Println(time.Now().Format(time.RFC822) + " : " + m)
}
