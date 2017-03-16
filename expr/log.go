package expr

import (
	"io"
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
	CRITICAL = g_logLevel.WrapFunc(int(LogCritical), log.Printf, "[expr][C] ")
	ERR      = g_logLevel.WrapFunc(int(LogError), log.Printf, "[expr][E] ")
	WARN     = g_logLevel.WrapFunc(int(LogWarning), log.Printf, "[expr][W] ")
	INFO     = g_logLevel.WrapFunc(int(LogInfo), log.Printf, "[expr][I] ")
	API      = g_logLevel.WrapFunc(int(LogApi), log.Printf, "[expr][A] ")
	DBG      = g_logLevel.WrapFunc(int(LogDebug), log.Printf, "[expr][D] ")
)

var g_logLevel logfn.LogLevel

func init() {
	log.SetFlags(log.Lmicroseconds | log.LstdFlags)
}

func SetLogLevel(n LogLevel) {
	g_logLevel.SetLevel(int(n))
}

func SetLogOutput(w io.Writer) {
	log.SetOutput(w)
}
