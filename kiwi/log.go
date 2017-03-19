package kiwi

import (
	"log"

	"github.com/noypi/logfn"
)

type LogLevel int

const (
	LogCritical LogLevel = iota
	LogError
	LogWarning
	LogInfo
	LogApi
	LogDebug
)

var (
	CRITICAL = g_logLevel.WrapFunc(int(LogCritical), log.Printf, "[kiwi][C] ")
	ERR      = g_logLevel.WrapFunc(int(LogError), log.Printf, "[kiwi][E] ")
	WARN     = g_logLevel.WrapFunc(int(LogWarning), log.Printf, "[kiwi][W] ")
	INFO     = g_logLevel.WrapFunc(int(LogInfo), log.Printf, "[kiwi][I] ")
	API      = g_logLevel.WrapFunc(int(LogApi), log.Printf, "[kiwi][A] ")
	DBG      = g_logLevel.WrapFunc(int(LogDebug), log.Printf, "[kiwi][D] ")
)

var g_logLevel logfn.LogLevel

func SetLogLevel(n LogLevel) {
	g_logLevel.SetLevel(int(n))
}
