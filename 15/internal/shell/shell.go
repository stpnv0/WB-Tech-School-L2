package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"shell/internal/builtins"
	"strings"
	"sync"
)

type Shell struct {
	reader *bufio.Reader
}

func New() *Shell {
	return &Shell{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (s *Shell) Run() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for {
			<-sigChan
			fmt.Println()
			fmt.Print("$ ")
		}
	}()

	for {
		fmt.Print("$ ")
		input, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("exit")
				return
			}
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if err := s.executeInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (s *Shell) executeInput(input string) error {
	pipeParts := strings.Split(input, "|")
	if len(pipeParts) == 0 {
		return nil
	}

	if len(pipeParts) == 1 {
		return s.runSingleCommand(pipeParts[0], os.Stdin, os.Stdout)
	}

	return s.executePipeline(pipeParts)
}

func (s *Shell) runSingleCommand(cmdStr string, stdin io.Reader, stdout io.Writer) error {
	args := strings.Fields(cmdStr)
	if len(args) == 0 {
		return nil
	}
	cmdName := args[0]
	cmdArgs := args[1:]

	if cmdName == "exit" {
		os.Exit(0)
	}

	if fn, ok := builtins.Builtins[cmdName]; ok {
		return fn(cmdArgs, stdin, stdout)
	}

	return s.runExternal(cmdName, cmdArgs, stdin, stdout)
}

func (s *Shell) runExternal(name string, args []string, stdin io.Reader, stdout io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func (s *Shell) executePipeline(parts []string) error {
	var wg sync.WaitGroup

	type proc struct {
		cmdStr string
		stdin  io.Reader
		stdout io.Writer
		closer io.Closer
	}

	var procs []proc

	var nextStdin io.Reader = os.Stdin

	for i, part := range parts {
		isLast := i == len(parts)-1
		currentCmd := strings.TrimSpace(part)

		var stdout io.Writer = os.Stdout
		var closer io.Closer

		if !isLast {
			r, w := io.Pipe()
			stdout = w
			closer = w
			procs = append(procs, proc{
				cmdStr: currentCmd,
				stdin:  nextStdin,
				stdout: stdout,
				closer: closer,
			})
			nextStdin = r
		} else {
			procs = append(procs, proc{
				cmdStr: currentCmd,
				stdin:  nextStdin,
				stdout: stdout,
				closer: nil,
			})
		}
	}

	for _, p := range procs {
		p := p
		wg.Add(1)

		// Parse cmd
		args := strings.Fields(p.cmdStr)
		if len(args) == 0 {
			if p.closer != nil {
				p.closer.Close()
			}
			wg.Done()
			continue
		}

		cmdName := args[0]
		cmdArgs := args[1:]

		go func() {
			defer wg.Done()
			defer func() {
				if p.closer != nil {
					p.closer.Close()
				}
			}()

			if fn, ok := builtins.Builtins[cmdName]; ok {
				// Builtin
				// CD in pipeline is ignored/doesn't affect parent (and here we are in goroutine anyway)
				// We run it
				if err := fn(cmdArgs, p.stdin, p.stdout); err != nil {
					fmt.Fprintf(os.Stderr, "builtin '%s' error: %v\n", cmdName, err)
				}
			} else {
				// External
				// We use exec.Command
				cmd := exec.Command(cmdName, cmdArgs...)
				cmd.Stdin = p.stdin
				cmd.Stdout = p.stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					// We just print error?
					// In pipeline, if one fails, what happens? "set -o pipefail"?
					// Basic shell: execution continues potentially?
					// Broken pipe might occur.
					fmt.Fprintf(os.Stderr, "command '%s' error: %v\n", cmdName, err)
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
