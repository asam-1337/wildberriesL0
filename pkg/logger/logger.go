package logger

import (
	"log"
	"os"
)

func Error() {
	l := log.New(os.Stdout, "[ERROR]")
}
