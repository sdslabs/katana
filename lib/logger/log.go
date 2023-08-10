package logger

import (
	"os"

	logging "github.com/op/go-logging"
)

const (
	ErrorTAG = "[ERROR]"
	InfoTAG  = "[INFO]"
)

var (
	log      = logging.MustGetLogger("katana")
	format   = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	LogInfo  = log.Info
	LogError = log.Error
)

// func getTimeStamp() string {
// 	return fmt.Sprintf("%v", time.Now().Unix())
// }

func init() {
	var infoBackend logging.Backend
	var errBackend logging.Backend

	infoBackend = logging.NewLogBackend(os.Stdout, InfoTAG, 0)
	errBackend = logging.NewLogBackend(os.Stderr, ErrorTAG, 0)

	infoBackend = logging.NewBackendFormatter(infoBackend, format)
	errBackend = logging.NewBackendFormatter(errBackend, format)

	errBackendLeveled := logging.AddModuleLevel(errBackend)
	errBackendLeveled.SetLevel(logging.ERROR, "")

	logging.SetBackend(infoBackend, errBackendLeveled)
}
