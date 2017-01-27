package logger

import (
	"errors"
	"sync"
	"time"
)

type LogFunc func(string, ...interface{})

type StateLogger interface {
	LogError(err error)
	Fixed()
	WithInterval(logErrorInterval time.Duration)
	WithSuccessLogger(logFunc LogFunc)
}

// NewStateLogger creates a new Logger that logs err only if logErrorInterval
// have passed from the last err, or it is a different err than the last seen.
func NewStateLogger(name string, logFunc LogFunc) StateLogger {
	return &stateLogger{
		name:             name,
		logError:         logFunc,
		logSuccess:       logFunc,
		logErrorInterval: logOnlyChanges,

		err:   errNoError,
		mutex: &sync.Mutex{},
	}
}

// Set logging interval, if time has passed since the last error
// and the error reoccur, it will be logged.
func (sl *stateLogger) WithInterval(logErrorInterval time.Duration) {
	sl.logErrorInterval = logErrorInterval
}

// Change the logging function for success
func (sl *stateLogger) WithSuccessLogger(logFunc LogFunc) {
	sl.logSuccess = logFunc
}

type stateLogger struct {
	name             string
	logError         LogFunc
	logSuccess       LogFunc
	logErrorInterval time.Duration

	err           error
	errLastLogged time.Time
	mutex         *sync.Mutex
}

const logOnlyChanges = time.Duration(-1)

var errNoError = errors.New("not an err")

// LogError logs an err if it is different from the last seen err,
// or that logErrorInterval have passed since the last reported err.
func (sl *stateLogger) LogError(err error) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()
	if !sl.shouldLogError(err) {
		return
	}
	sl.err = err
	sl.errLastLogged = time.Now()

	sl.logError("[%s] error: %s", sl.name, sl.err)
}

// Fixed makes the stateLogger understand that the state is fixed, and when
// the next err will occur, it will log it.
func (sl *stateLogger) Fixed() {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()
	if sl.logErrorInterval == 0 || sl.err == errNoError {
		return
	}
	sl.err = errNoError

	sl.logSuccess("[%s] fixed!", sl.name)
}

func (sl *stateLogger) shouldLogError(err error) bool {
	timeCause := sl.logErrorInterval != logOnlyChanges && time.Since(sl.errLastLogged) > sl.logErrorInterval
	return timeCause || err.Error() != sl.err.Error()
}
