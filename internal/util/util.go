package util

import (
	"fmt"
	"os"
)

func Must(err error) {
	if err != nil {
		fmt.Printf("fatal: %s\n", err.Error())
		os.Exit(-1)
	}
}
