package builtins

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"os/exec"
)

type BuiltinFunc func(args []string, stdin io.Reader, stdout io.Writer) error

var Builtins = map[string]BuiltinFunc{
	"cd":   cd,
	"pwd":  pwd,
	"echo": echo,
	"kill": kill,
	"ps":   mn_ps,
}

func cd(args []string, stdin io.Reader, stdout io.Writer) error {
	dir := ""
	if len(args) > 0 {
		dir = args[0]
	} else {
		var err error
		dir, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	}
	return os.Chdir(dir)
}

func pwd(args []string, stdin io.Reader, stdout io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, dir)
	return nil
}

func echo(args []string, stdin io.Reader, stdout io.Writer) error {
	fmt.Fprintln(stdout, strings.Join(args, " "))
	return nil
}

func kill(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: kill <pid>")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid pid: %s", args[0])
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}

func mn_ps(args []string, stdin io.Reader, stdout io.Writer) error {
	cmd := exec.Command("ps", "-A")
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
