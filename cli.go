package main

import "io"

const (
	ExitCodeOk             = 0
	ExitCodeParseFlagError = 1
	ExitCodeError          = 2
)

type CLI struct {
	outStream, errStream io.Writer
}

func (c *CLI) Run(args []string) int {
	return ExitCodeParseFlagError
}
