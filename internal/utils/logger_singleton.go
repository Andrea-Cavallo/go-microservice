package utils

import (
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	once   sync.Once
	logger *logrus.Logger
)

// GetLogger ritorna l'istanza singleton del logger
func GetLogger() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FullTimestamp:   true,
		})
		logger.SetReportCaller(true)

		// Configura l'output del logger per scrivere sia su console che su file
		logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.Fatalf("Failed to open log file: %v", err)
		}

		// Imposta l'output del logger su console e su file
		logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	})
	return logger
}

// WithContext aggiunge informazioni contestuali all'entry di log
func WithContext() *logrus.Entry {
	return GetLogger().WithField("caller", getCallerInfo())
}

// getCallerInfo recupera il nome del file e la linea da cui viene chiamata la funzione
func getCallerInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	return fn.Name() + ":" + file + ":" + string(rune(line))
}
