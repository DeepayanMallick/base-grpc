package logging

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	isTerminal = terminal.IsTerminal(int(os.Stdout.Fd()))
	isTest     = strings.HasSuffix(os.Args[0], ".test")
	flagMu     = sync.Mutex{}
)

func NewLogger(c *viper.Viper) *logrus.Logger {
	logger := logrus.New()

	switch c.GetString("runtime.loglevel") {
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	if isTest {
		testing.Init()

		flagMu.Lock()
		if !flag.Parsed() {
			flag.Parse()
		}
		flagMu.Unlock()

		if !testing.Verbose() {
			logger.Level = logrus.FatalLevel
		}
	}

	logger.Out = os.Stdout
	switch {
	case c != nil && c.GetString("runtime.logformat") == "json":
		logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     true,
		}
	case isTerminal:
		logger.Formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		}
	default:
		logger.Formatter = &FluentdFormatter{
			TimestampFormat: time.RFC3339,
		}
	}

	logger.SetReportCaller(true)
	return logger
}

func WithError(err error, logger logrus.FieldLogger) *logrus.Entry {
	return logger.WithError(err).WithField("stacktrace", fmt.Sprintf("%+v", err))
}
