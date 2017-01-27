package logger

import (
	"log"
	"os"
	"errors"
	"time"
)

// ExampleStateLoggerWhy shows an example for a logging that you want to
// use a state logger for
func ExampleStateLoggerWhy() {
	l := log.New(os.Stdout, "", 0)

	// It is common to handle an err by logging it
	do := fail4times()
	for {
		err := do()
		if err != nil {
			l.Println(err)
			continue
		}
		break
	}

	// Your logs might get really massey:

	// Output: fails
	// fails
	// fails
	// fails
}

// ExampleStateLoggerHow shows an example who to use the StateLogger and why
// it is good for.
func ExampleStateLoggerHow() {
	l := log.New(os.Stdout, "", 0)

	// You could use the state logger for this:

	do := fail4times()
	sl := NewStateLogger("do state", l.Printf, l.Printf, time.Minute)

	for {
		err := do()
		if err != nil {
			sl.LogError(err)
			continue
		}
		sl.Fixed()
		break
	}

	// Your logs might get really massey:

	// Output: [do state] error: fails
	// [do state] fixed!
}

// fail4times returns a function that fails the 4 first time it was called
func fail4times() (func() error) {
	i := 0
	return func() error {
		if i >= 4 {
			return nil
		}
		i += 1
		return errors.New("fails")
	}
}

