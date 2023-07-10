package utils

import (
	"time"
)

var Ticker *time.Ticker

func InitTicker(duration time.Duration) {
	Ticker = time.NewTicker(duration)
}

func GetTicker() *time.Ticker {
	return Ticker
}

func SetTicker(ticker *time.Ticker) {
	Ticker = ticker
}

func StopTicker() {
	Ticker.Stop()
}

func ResetTicker(duration time.Duration) {
	Ticker.Reset(duration)
}

func GetRemainingTimeBeforeNextTick() time.Duration {
	return time.Until(<-Ticker.C)
}
