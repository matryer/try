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
	try.MaxRetries = 20
	SomeFunction := func() (string, error) {
		return "", nil
	}
	var value string
	err := try.Do(func(attempt int) (error, bool) {
		var err error
		value, err = SomeFunction()
		return err, attempt < 5 // try 5 times
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
	err := try.Do(func(attempt int) (err error, retry bool) {
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
	err := try.Do(func(attempt int) (error, bool) {
		callCount++
		return nil, attempt < 5
	})
	is.NoErr(err)
	is.Equal(callCount, 1)
}

func TestTryDoFailed(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (error, bool) {
		callCount++
		return theErr, attempt < 5
	})
	is.Equal(err, theErr)
	is.Equal(callCount, 5)
}

func TestTryPanics(t *testing.T) {
	is := is.New(t)
	theErr := errors.New("something went wrong")
	callCount := 0
	err := try.Do(func(attempt int) (err error, retry bool) {
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
	is.Equal(err.Error(), "panic: I don't like three")
	is.Equal(callCount, 5)
}

func TestRetryLimit(t *testing.T) {
	is := is.New(t)
	err := try.Do(func(attempt int) (error, bool) {
		return errors.New("nope"), true
	})
	is.OK(err)
	is.Equal(try.IsMaxRetries(err), true)
}
