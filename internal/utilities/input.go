package utilities

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func ReadValue(value string, echo bool, required bool) (string, error) {
	os.Stdout.WriteString(value + ": ")

	var byteValue []byte
	var err error

	if echo {
		byteValue, err = term.ReadPassword(int(syscall.Stdin))
	} else {
		byteValue, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
		byteValue = byteValue[:len(byteValue)-1]
	}
	if err != nil {
		return "", err
	}

	os.Stdout.WriteString("\n")

	if required && string(byteValue) == "" {
		return "", fmt.Errorf("%s is required", value)
	}
	return string(byteValue), nil
}
