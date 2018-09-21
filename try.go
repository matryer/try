package try

import (
	"errors"
	"fmt"
)

// Func represents functions that can be retried.
type Func func(attempt int) (retry bool, err error)

// Do keeps trying the function until the second argument
// returns false, or no error is returned.
func Do(fn Func) error {
	var err error
	var retry bool
	attempt := 1
	for {
		retry, err = fn(attempt)
		if err == nil {
			break
		}
		attempt++
		if !retry {
			return errors.New(fmt.Sprintf("exceeded retry limit - %v", err))
		}
	}
	return err
}
