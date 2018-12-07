// +build !appengine

package log

import (
	"io"

	colorable "github.com/mattn/go-colorable"
)

func output() io.Writer {
	return colorable.NewColorableStdout()
}
