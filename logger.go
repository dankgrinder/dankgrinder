// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// logFileWriter is an io.Writer that outputs to a file for logging.
type logFileWriter struct {
	username string
	cluster  string
	dir      string
}

// logFileHook is a logrus hook that writes all logs on a logger to a file as
// well.
type logFileHook struct {
	dir string
}

// stdLoggerHook is a logrus hook that puts some of the logs on a logger on the
// standard logger as well.
type stdLoggerHook struct {
	cluster  string
	username string
	verbose  bool
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

	cld := path.Join(ld, lfw.cluster)
	if _, err := os.Stat(cld); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(cld, 0755); err != nil {
				return 0, fmt.Errorf("error while creating cluster logs dir: %v", err)
			}
		}
	}

	uld := path.Join(cld, lfw.username)
	if _, err := os.Stat(uld); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(uld, 0755); err != nil {
				return 0, fmt.Errorf("error while creating cluster intance log dir: %v", err)
			}
		}
	}

	date := time.Now().Format("02-01-2006")

	name := path.Join(uld, fmt.Sprintf("%v.log", date))
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
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
	levels := []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
	if slh.verbose {
		levels = append(levels, []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		}...)
	}
	return levels
}

func (slh stdLoggerHook) Fire(e *logrus.Entry) error {
	fields := e.Data
	fields["instance"] = slh.username
	fields["cluster"] = slh.cluster
	logrus.WithFields(fields).Log(e.Level, e.Message)
	return nil
}

type instanceLoggerOpts struct {
	username             string
	discriminator        string
	cluster              string
	dir                  string
	id                   string
	debug                bool
	verboseStdLoggerHook bool
}

func newInstanceLogger(opts instanceLoggerOpts) *logrus.Logger {
	rawUsername, cleanedUsername := opts.username, ""
	allowedChars := regexp.MustCompile(`[a-zA-Z0-9-_]`)
	rawUsername = strings.Replace(rawUsername, " ", "-", -1)
	rawUsername = strings.ToLower(rawUsername)

	for _, char := range rawUsername {
		if allowedChars.MatchString(string(char)) {
			cleanedUsername += string(char)
		}
	}

	if cleanedUsername == "" {
		cleanedUsername = opts.id
	}

	logger := logrus.New()
	if opts.debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger = logrus.New()
	logger.SetOutput(ioutil.Discard)
	if opts.dir != "" {
		logger.SetOutput(logFileWriter{
			username: fmt.Sprintf("%v#%v", cleanedUsername, opts.discriminator),
			cluster:  opts.cluster,
			dir:      opts.dir,
		})
	}
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.AddHook(stdLoggerHook{
		username: fmt.Sprintf("%v#%v", opts.username, opts.discriminator),
		cluster:  opts.cluster,
		verbose:  opts.verboseStdLoggerHook,
	})
	return logger
}
