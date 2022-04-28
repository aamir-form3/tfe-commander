package log

import (
	"io"
	"os"
)

var Writer io.Writer = os.Stdout
