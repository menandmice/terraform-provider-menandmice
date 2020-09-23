package diag

import "fmt"

type Diagnostics error

func FromErr(err error) error {
	return err
}

func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
func Append(err1, err2 Diagnostics) Diagnostics {
	return err1
	// return  append(err1, err2..)
}
