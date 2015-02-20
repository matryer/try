// Package try provides retry functionality.
//     var value string
//     err := try.Do(func(attempt int) (error, bool) {
//       var err error
//       value, err = SomeFunction()
//       return err, attempt < 5 // try 5 times
//     })
//     if err != nil {
//       log.Fatalln("error:", err)
//     }
package try
