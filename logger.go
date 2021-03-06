// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

type logFileWriter struct {
	username string
	dir      string
}

type logFileHook struct {
	dir string
}

type stdLoggerHook struct {
	username string
}

func (lfw logFileWriter) Write(b []byte) (int, error) {
	ld := path.Join(lfw.dir, "logs")
	if _, err := os.Stat(ld); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(ld, 0755); err != nil {
				return 0, fmt.Errorf("error while creating logs dir: %v", err)
			}
		}
	}
	uld := path.Join(lfw.dir, fmt.Sprintf("logs/%v", lfw.username))
	if _, err := os.Stat(uld); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(uld, 0755); err != nil {
				return 0, fmt.Errorf("error while creating logs dir: %v", err)
			}
		}
	}
	date := time.Now().Format("02-01-2006")
	name := fmt.Sprintf("logs/%v/%v.log", lfw.username, date)
	f, err := os.OpenFile(path.Join(lfw.dir, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
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
	ld := path.Join(lfh.dir, "logs")
	if _, err := os.Stat(ld); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(ld, 0755); err != nil {
				return fmt.Errorf("error while creating logs dir: %v", err)
			}
		}
	}
	date := time.Now().Format("02-01-2006")
	name := fmt.Sprintf("logs/dankgrinder-%v.log", date)
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
	fields := e.Data
	fields["instance"] = slh.username
	logrus.WithFields(fields).Log(e.Level, e.Message)
	return nil
}

func newInstanceLogger(username, dir string, debug bool) *logrus.Logger {
	logger := logrus.New()
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger = logrus.New()
	logger.SetOutput(logFileWriter{
		username: username,
		dir:      dir,
	})
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.AddHook(stdLoggerHook{username: username})
	return logger
}
