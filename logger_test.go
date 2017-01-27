package logger

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

const (
	name     = "state"
	safeWait = 100 * time.Millisecond
)

type mockLogger struct {
	mock.Mock
}

func (l *mockLogger) err(msg string, args ...interface{}) {
	callArgs := append(append([]interface{}{}, msg), args...)
	l.Called(callArgs...)
}

func (l *mockLogger) fixed(msg string, args ...interface{}) {
	callArgs := append(append([]interface{}{}, msg), args...)
	l.Called(callArgs...)
}

func TestStateLogger(t *testing.T) {
	err1 := errors.New("err 1")
	err2 := errors.New("err 2")
	name := "state"

	log := new(mockLogger)
	sl := NewStateLogger(name, log.err)
	sl.WithInterval(safeWait)
	sl.WithSuccessLogger(log.fixed)

	log.On("err", mock.Anything, name, err1).Return(nil).Once()
	sl.LogError(err1)
	log.AssertNumberOfCalls(t, "err", 1)
	log.AssertNumberOfCalls(t, "fixed", 0)

	log.On("err", mock.Anything, name, err2).Return(nil).Once()
	sl.LogError(err2)
	log.AssertNumberOfCalls(t, "err", 2)
	log.AssertNumberOfCalls(t, "fixed", 0)

	log.On("err", mock.Anything, name, err1).Return(nil).Once()
	sl.LogError(err1)
	sl.LogError(err1)
	log.AssertNumberOfCalls(t, "err", 3)
	log.AssertNumberOfCalls(t, "fixed", 0)

	time.Sleep(safeWait)

	log.On("err", mock.Anything, name, err1).Return(nil).Once()
	sl.LogError(err1)
	sl.LogError(err1)
	log.AssertNumberOfCalls(t, "err", 4)
	log.AssertNumberOfCalls(t, "fixed", 0)

	log.On("err", mock.Anything, name, err2).Return(nil).Once()
	sl.LogError(err2)
	log.AssertNumberOfCalls(t, "err", 5)
	log.AssertNumberOfCalls(t, "fixed", 0)

	log.On("fixed", mock.Anything, name).Return(nil).Once()
	sl.Fixed()
	log.AssertNumberOfCalls(t, "err", 5)
	log.AssertNumberOfCalls(t, "fixed", 1)

	sl.Fixed()
	log.AssertNumberOfCalls(t, "err", 5)
	log.AssertNumberOfCalls(t, "fixed", 1)

	log.On("err", mock.Anything, name, err2).Return(nil).Once()
	sl.LogError(err2)
	log.AssertNumberOfCalls(t, "err", 6)
	log.AssertNumberOfCalls(t, "fixed", 1)

	log.On("fixed", mock.Anything, name).Return(nil).Once()
	sl.Fixed()
	log.AssertNumberOfCalls(t, "err", 6)
	log.AssertNumberOfCalls(t, "fixed", 2)
}

func TestStateLoggerAlwaysLog(t *testing.T) {
	err1 := errors.New("err 1")
	err2 := errors.New("err 2")

	m := new(mockLogger)
	sl := NewStateLogger(name, m.err)
	sl.WithInterval(0)
	sl.WithSuccessLogger(m.fixed)

	m.On("err", mock.Anything, name, err1).Return(nil).Once()
	sl.LogError(err1)
	m.AssertNumberOfCalls(t, "err", 1)
	m.AssertNumberOfCalls(t, "fixed", 0)

	m.On("err", mock.Anything, name, err1).Return(nil).Once()
	sl.LogError(err1)
	m.AssertNumberOfCalls(t, "err", 2)
	m.AssertNumberOfCalls(t, "fixed", 0)

	m.On("err", mock.Anything, name, err2).Return(nil).Once()
	sl.LogError(err2)
	m.AssertNumberOfCalls(t, "err", 3)
	m.AssertNumberOfCalls(t, "fixed", 0)

	m.On("err", mock.Anything, name, err2).Return(nil).Once()
	sl.LogError(err2)
	m.AssertNumberOfCalls(t, "err", 4)
	m.AssertNumberOfCalls(t, "fixed", 0)

	m.On("err", mock.Anything, name, err1).Return(nil).Once()
	sl.LogError(err1)
	m.AssertNumberOfCalls(t, "err", 5)
	m.AssertNumberOfCalls(t, "fixed", 0)

	sl.Fixed()
	m.AssertNumberOfCalls(t, "err", 5)
	m.AssertNumberOfCalls(t, "fixed", 0)

	sl.Fixed()
	m.AssertNumberOfCalls(t, "err", 5)
	m.AssertNumberOfCalls(t, "fixed", 0)

	m.On("err", mock.Anything, name, err2).Return(nil).Once()
	sl.LogError(err2)
	m.AssertNumberOfCalls(t, "err", 6)
	m.AssertNumberOfCalls(t, "fixed", 0)
}

func TestStateFirstFixed(t *testing.T) {
	m := new(mockLogger)
	sl := NewStateLogger("state", m.err)
	sl.WithSuccessLogger(m.fixed)

	sl.Fixed()
	m.AssertNumberOfCalls(t, "err", 0)
	m.AssertNumberOfCalls(t, "fixed", 0)

	sl = NewStateLogger("state", m.err)
	sl.WithSuccessLogger(m.fixed)
	sl.WithInterval(0)

	sl.Fixed()
	m.AssertNumberOfCalls(t, "err", 0)
	m.AssertNumberOfCalls(t, "fixed", 0)
}

func TestStateErrorsWithTheSameMessage(t *testing.T) {
	err := errors.New("err 1")
	errCopy := errors.New("err 1")

	log := new(mockLogger)
	sl := NewStateLogger("state", log.err)
	sl.WithSuccessLogger(log.fixed)

	log.On("err", mock.Anything, name, err).Return(nil).Once()
	sl.LogError(err)
	log.AssertNumberOfCalls(t, "err", 1)
	log.AssertNumberOfCalls(t, "fixed", 0)

	sl.LogError(errCopy)
	log.AssertNumberOfCalls(t, "err", 1)
	log.AssertNumberOfCalls(t, "fixed", 0)
}
