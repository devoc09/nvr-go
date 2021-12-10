package main

import "os"

func main() {
	c := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(c.Run(os.Args))
}
