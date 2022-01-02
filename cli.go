package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
)

const (
	ExitCodeOk             = 0
	ExitCodeParseFlagError = 1
	ExitCodeError          = 2
)

type CLI struct {
	outStream, errStream io.Writer
}

func (c *CLI) Run(args []string) int {
	var remote, remotewait bool

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.Usage = func() {
		fmt.Fprintf(c.errStream, usage, Name)
	}
	flags.BoolVar(&remote, "r", false, "Execute :edit to open file")
	flags.BoolVar(&remote, "rw", false, "Execute :edit to open file")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if remote {
		nv, err := nvim.Dial(os.Getenv("NVIM_LISTEN_ADDRESS"), nvim.DialServe(true))
		// nv, err := nvim.Dial(os.Getenv("NVIM_LISTEN_ADDRESS"), nvim.DialServe(false))
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed dial up remote neovim.")
			return ExitCodeError
		}
		defer nv.Close()

		var cmd []string
		if len(args) == 2 {
			cmd = append(cmd, "edit")
		} else {
			cmd = append(cmd, "edit", args[2])
		}
		if err := nv.Command(strings.Join(cmd, " ")); err != nil {
			fmt.Fprintf(c.errStream, "Error Exec neovim command :edit")
			return ExitCodeError
		}
	}

	if remotewait {
		nv, err := nvim.Dial(os.Getenv("NVI_LISTEN_ADDRESS"), nvim.DialServe(true))
		if err != nil {
			fmt.Fprintf(c.errStream, "Failed dial up remote neovim.")
			return ExitCodeError
		}
		defer nv.Close()

		var cmd []string
		if len(args) == 2 {
			cmd = append(cmd, "edit")
		} else {
			cmd = append(cmd, "edit", args[2])
		}

		if err = waitCurrentBuf(c, nv); err != nil {
			fmt.Fprintf(c.errStream, "Error Wait Current Buffer.")
			return ExitCodeError
		}

		if err := nv.Command(strings.Join(cmd, " ")); err != nil {
			fmt.Fprintf(c.errStream, "Error Exec neovim comamnd :edit")
			return ExitCodeError
		}
	}

	return ExitCodeOk
}

// I haven't used it yet.
// func isPipeInput(fd int) bool {
// 	return !term.IsTerminal(fd)
// }

func waitCurrentBuf(c *CLI, nvim *nvim.Nvim) error {
	chanid := nvim.ChannelID()

	cmds := []string{
		"augroup nvr-go",
		"autocmd BufDelete <buffer> silent! call rpcnotify(%d, 'BufDelete')",
		"autocmd VimLeave * if exists('v:exiting') && v:exiting > 0 | silent! call rpcnotify(%d, 'Exit', v:exiting) | endif",
		"augroup END",
	}

	for i, cmd := range cmds {
		if i == 1 || i == 2 {
			if err := nvim.Command(fmt.Sprintf(cmd, chanid)); err != nil {
				fmt.Fprintf(c.errStream, "Error Exec neovim command: %s", cmd)
				return err
			}
		} else {
			if err := nvim.Command(cmd); err != nil {
				fmt.Fprintf(c.errStream, "Error Exec neovim command: %s", cmd)
				return err
			}
		}
	}
	return nil
}

const usage = `
Usage: %s [options]
`
