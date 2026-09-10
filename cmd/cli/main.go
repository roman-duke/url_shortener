package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/roman-duke/url_shortener/shortener"
	"github.com/roman-duke/url_shortener/storage/memory"
)

var (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

var ErrSvrError = errors.New("Internal Server Error")

func main() {
	// create the store
	store := memory.MyStore{}

	svc := shortener.NewService(&store)

	// print the welcome message
	fmt.Printf("\n%s============================================================================\n", Cyan)
	fmt.Print(`
	█  █ ███  █        ████ █  █ ████ ███  ████ ████ █  █ ████ ███
	█  █ █  █ █        █    █  █ █  █ █  █  █   █    ██ █ █    █  █
	█  █ ███  █        ████ ████ █  █ ███   █   ███  █ ██ ███  ███
	█  █ █ █  █           █ █  █ █  █ █ █   █   █    █  █ █    █ █
	████ █  █ ████     ████ █  █ ████ █  █  █   ████ █  █ ████ █  █
	`)
	fmt.Print("\n")
	fmt.Printf("============================================================================\n\n%s", Reset)

	// Parse the command line arguments
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// use print so that the user's actual commands stay on the same line
		fmt.Print("shortener> ")

		if !scanner.Scan() {
			break
		}

		// now parse the user's input to get the command and the args
		tokens := strings.Fields(scanner.Text())

		if len(tokens) == 0 {
			continue
		} else if len(tokens) < 2 {
			fmt.Fprintf(os.Stderr, "%scommand needs at least one argument%s\n", Red, Reset)
			continue
		}

		cmd := tokens[0]
		args := tokens[1:]

		dispatch(svc, cmd, args, os.Stdout, os.Stderr)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
