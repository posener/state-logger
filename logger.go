package logger

import (
	"fmt"
	"sync"
	"time"
	"errors"
)

var errNoError = errors.New("not an err")

type LogFunc func(string, ...interface{})

type StateLogger interface {
	LogError(err error)
	Fixed()
}

// NewStateLogger creates a new Logger that logs err only if logErrorInterval
// have passed from the last err, or it is a different err than the last seen.
func NewStateLogger(name string, logError LogFunc, logSuccess LogFunc, logErrorInterval time.Duration) StateLogger {
	return &stateLogger{
		name: name,
		logError:         logError,
		logSuccess:       logSuccess,
		logErrorInterval: logErrorInterval,

		err:   errNoError,
		mutex: &sync.Mutex{},
	}
}

type stateLogger struct {
	name string
	logError         LogFunc
	logSuccess       LogFunc
	logErrorInterval time.Duration

	err           error
	errLastLogged time.Time
	mutex         *sync.Mutex
}

// LogError logs an err if it is different from the last seen err,
// or that logErrorInterval have passed since the last reported err.
func (sl *stateLogger) LogError(err error) {
	sl.mutex.Lock()
	if err.Error() == sl.err.Error() && time.Since(sl.errLastLogged) < sl.logErrorInterval {
		sl.mutex.Unlock()
		return
	}
	sl.err = err
	sl.errLastLogged = time.Now()
	msg := sl.formatMessage()
	sl.mutex.Unlock()

	sl.logError(msg)
}

// Fixed makes the stateLogger understand that the state is fixed, and when
// the next err will occur, it will log it.
func (sl *stateLogger) Fixed() {
	sl.mutex.Lock()
	if sl.logErrorInterval == 0 || sl.err == errNoError {
		sl.mutex.Unlock()
		return
	}
	sl.err = errNoError
	msg := sl.formatMessage()
	sl.mutex.Unlock()

	sl.logSuccess(msg)
}

func (sl *stateLogger) formatMessage() string {
	if sl.err != errNoError {
		return fmt.Sprintf("[%s] error: %s ", sl.name, sl.err)
	} else {
		return fmt.Sprintf("[%s] fixed!", sl.name)
	}
}
