package main

import (
	"errors"
	"fmt"
	"os"
)

var Error100 = &CustomError{msg: "não foi possível executar o processo", code: 100}

type CustomError struct {
	msg  string
	code int
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("%s %d", c.msg, c.code)
}

func NewCustomError(msg string, code int) error {
	return &CustomError{msg: msg, code: code}
}

func process() (string, error) {
	f, err := os.Open("foo")
	if err != nil {
		return "", fmt.Errorf("não foi possíve executar o processo: (%w)", err)
	}
	return f.Name(), nil
}

func main() {
	r, err := process()
	if err != nil {
		// var err100 *CustomError
		// if errors.As(err, &err100) {
		// 	fmt.Println(err100.code)
		// }
		// return
		fmt.Println(err)
		fmt.Println(errors.Unwrap(err))
	}
	fmt.Println(r)
}
