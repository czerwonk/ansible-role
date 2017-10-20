package main

import (
	"fmt"
	"io"
)

type debugWriter struct {
	writer io.Writer
}

func (d *debugWriter) Write(p []byte) (n int, err error) {
	fmt.Print(string(p))
	return d.writer.Write(p)
}
