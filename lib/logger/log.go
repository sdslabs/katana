package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/sdslabs/katana/configs"
)

const (
	ErrorTAG = 1
	InfoTAG  = 2
	DebugTAG = 3
)

var tagToString = map[int]string{
	ErrorTAG: "[ERROR]",
	InfoTAG:  "[INFO]",
	DebugTAG: "[DEBUG]",
}

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var logFile *os.File
var fileLogging bool

func getTimeStamp() string {
	return fmt.Sprintf("%v", time.Now().Unix())
}

func Log(tag int, messages ...string) {
	fmt.Fprintf(logFile, "abc")
}

func LogError(message string) {
	fmt.Println()
}

func init() {
	var err error
	logFile, err = os.Open(configs.KatanaConfig.LogFile)
	if err != nil {
		LogError("Log file could not be accessed")
	} else {
		fileLogging = true
	}
}
