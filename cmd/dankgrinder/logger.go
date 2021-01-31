package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

type fileLogger struct {
	username string
	dir      string
}

type logFileHook struct {
	dir string
}

type stdLoggerHook struct {}

func (fl fileLogger) Write(b []byte) (int, error) {
	date := time.Now().Format("02-01-2006")
	name := fmt.Sprintf("%v-%v.log", fl.username, date)
	f, err := os.OpenFile(path.Join(fl.dir, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(b)
}

func (lfh logFileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (lfh logFileHook) Fire(e *logrus.Entry) error {
	date := time.Now().Format("02-01-2006")
	name := fmt.Sprintf("dankgrinder-%v.log", date)
	f, err := os.OpenFile(path.Join(lfh.dir, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := (&logrus.JSONFormatter{}).Format(e)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (slh stdLoggerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (slh stdLoggerHook) Fire(e *logrus.Entry) error {
	logrus.WithFields(e.Data).Log(e.Level, e.Message)
	return nil
}

func newInstanceLogger(username, dir string) *logrus.Logger {
	logger := logrus.New()
	if cfg.Features.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger = logrus.New()
	logger.SetOutput(fileLogger{
		username: username,
		dir:      dir,
	})
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.AddHook(stdLoggerHook{})
	return logger
}
