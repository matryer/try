package try_test

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/try"
)

func TestTryExample(t *testing.T) {
	SomeFunction := func() (string, error) {
		return "", nil
	}
	var value string
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		value, err = SomeFunction()
		return attempt < 5, err // try 5 times
	})
	if err != nil {
		log.Fatalln("error:", err)
	}
}

func TestTryExamplePanic(t *testing.T) {
	SomeFunction := func() (string, error) {
		panic("something went badly wrong")
	}
	var value string
	err := try.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5 // try 5 times
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()
		value, err = SomeFunction()
		return
	})
	if err != nil {
		//log.Fatalln("error:", err)
	}
}

func TestTryDoSuccessful(t *testing.T) {
	is := is.New(t)
	callCount := 0
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, nil
	})
	is.NoErr(err)
	is.Equal(callCount, 1)
}

func TestTryDoFailed(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (bool, error) {
		callCount++
		return attempt < 5, theErr
	})
	is.Equal(err.Error(), "exceeded retry limit - something went wrong")
	is.Equal(callCount, 5)
}

func TestTryPanics(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (retry bool, err error) {
		retry = attempt < 5
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("panic: %v", r))
			}
		}()
		callCount++
		if attempt > 2 {
			panic("I don't like three")
		}
		err = theErr
		return
	})
	is.Equal(err.Error(), "exceeded retry limit - panic: I don't like three")
	is.Equal(callCount, 5)
}

func TestRetryLimit(t *testing.T) {
	is := is.New(t)
	SomeFunction := func() (string, error) {
		return "", errors.New("something went wrong")
	}
	var lastAttempt int
	err := try.Do(func(attempt int) (retry bool, err error) {
		_, err = SomeFunction()
		retry = attempt < 15 // try 15 times
		lastAttempt = attempt
		return retry, err
	})

	is.Equal(err.Error(), "exceeded retry limit - something went wrong")
	is.Equal(lastAttempt, 15)
}
