package cmd

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitializeLogger creates a new logger for use with the cli app.
func InitializeLogger(vcfg *viper.Viper) *log.Logger {
	var formatter log.Formatter

	switch vcfg.GetString("formatter") {
	case "default", "text":
		formatter = &log.TextFormatter{FullTimestamp: true}
	case "json":
		formatter = &log.JSONFormatter{}
	default:
		formatter = &log.TextFormatter{FullTimestamp: true}
	}

	logLevel := log.InfoLevel
	if vcfg.GetBool("debug") {
		logLevel = log.DebugLevel
	}

	return &log.Logger{
		Out:       os.Stdout,
		Formatter: formatter,
		Hooks:     make(log.LevelHooks),
		Level:     logLevel,
	}
}

type callBackFunc func()

// TermListener waits on os Interrupts, SIGTERM, or SIGINTs.
// When a signal is received, the callbacks are all called in
// parallel and will wait until whatever cleanup operations have run.
type TermListener struct {
	Ch chan os.Signal
	wg sync.WaitGroup

	Callbacks []callBackFunc
}

// AddCallback Adds a callback onto the term listener.
func (tl *TermListener) AddCallback(f callBackFunc) {
	tl.wg.Add(1)
	tl.Callbacks = append(tl.Callbacks, f)
}

// WaitForCtrlC will listen for an Interrupt, SIGTERM, or SIGINT
// before calling all callbacks in parallel and waiting for completion.
func (tl *TermListener) WaitForCtrlC() {
	tl.Ch = make(chan os.Signal, 2)
	signal.Notify(tl.Ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-tl.Ch

	for _, cb := range tl.Callbacks {
		go func(cb callBackFunc) {
			cb()
			tl.wg.Done()
		}(cb)
	}
	tl.wg.Wait()
}
