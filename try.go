package try

import "errors"

// MaxRetries is the maximum number of retries before bailing.
var MaxRetries = 10

var errMaxRetriesReached = errors.New("exceeded retry limit")

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func Do(fn func(attempt int) (err error, retry bool)) error {
	var err error
	var cont bool
	attempt := 1
	for {
		err, cont = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > MaxRetries {
			return errMaxRetriesReached
		}
	}
	return err
}

// IsMaxRetries checks whether the error is due to hitting the
// maximum number of retries or not.
func IsMaxRetries(err error) bool {
	return err == errMaxRetriesReached
}
