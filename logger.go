// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

//base
import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

type fileLogger struct {
	username string
	dir      string
}

//64
type logFileHook struct {
	dir string
}

type stdLoggerHook struct {
	username string
}

//this might be hum.. log formating
func (fl fileLogger) Write(b []byte) (int, error) {
	if _, err := os.Stat("logs"); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir("logs", 0755); err != nil {
				return 0, fmt.Errorf("error while creating logs dir: %v", err)
			}
		}
	}
	date := time.Now().Format("02-01-2006")
	name := fmt.Sprintf("logs/%v-%v.log", fl.username, date)
	f, err := os.OpenFile(path.Join(fl.dir, name), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(b)
}

//if you dont know wat this is... (its error levels)
func (lfh logFileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel, //panic!!!
	}
}

//error for creating log... or dir i think
func (lfh logFileHook) Fire(e *logrus.Entry) error {
	if _, err := os.Stat("logs"); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir("logs", 0755); err != nil {
				return fmt.Errorf("error while creating logs dir: %v", err)
			}
		}
	}
	//...this creats the time prefix ... right?
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

//this portion is what shows up in console or somthin
func (slh stdLoggerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel, //paniccc!
	}
}

//idek
func (slh stdLoggerHook) Fire(e *logrus.Entry) error {
	fields := e.Data
	fields["instance"] = slh.username
	logrus.WithFields(fields).Log(e.Level, e.Message)
	return nil
}

//may still be prefix code idk anymore
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
	logger.AddHook(stdLoggerHook{username: username})
	return logger
}
