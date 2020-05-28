package fmt

import (
	"fmt"
	"strings"
)

// debug util for tracking down source of stdout messages
// find/replace "fmt" with "github.com/powersjcb/monitor/go/src/lib/fmt" and update target string

const targetString = "EOF"

func Println(a ...interface{}) (n int, err error) {
	if strings.Contains(fmt.Sprintln(a...), targetString) {
		panic(a)
	}
	return fmt.Println(a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	if strings.Contains(fmt.Sprintf(format, a...), targetString) {
		panic(a)
	}
	return fmt.Printf(format, a...)
}

func Sprintf(format string, a ...interface{}) string {
	if strings.Contains(fmt.Sprintf(format, a...), targetString) {
		panic(a)
	}
	return fmt.Sprintf(format, a...)
}

func Errorf(format string, a ...interface{}) error {
	if strings.Contains(fmt.Sprintf(format, a...), targetString) {
		panic(a)
	}
	return fmt.Errorf(format, a...)
}
