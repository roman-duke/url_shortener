package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/roman-duke/url_shortener/shortener"
)

func dispatch(svc *shortener.Service, cmd string, args []string, out io.Writer, errOut io.Writer) {
	// for now, we just take the take only the first argument and ignore the rest
	// ("for now", like I would ever come back to this implementation)
	// "know thyself!"
	arg := args[0]

	switch cmd {
	case "shorten":
		result, err := svc.Shorten(arg)

		if errors.Is(err, shortener.ErrInvalidUrl) {
			fmt.Fprintf(errOut, "%s%s%s\n", Red, "Url is invalid", Reset)
		} else if err != nil {
			fmt.Fprintf(errOut, "%s%s%s\n", Red, ErrSvrError.Error(), Reset)
		} else {
			fmt.Fprintf(out, "%s%s%s\n", Green, result.Code, Reset)
		}

	case "resolve":
		result, err := svc.Resolve(arg)

		if errors.Is(err, shortener.ErrNotFound) {
			fmt.Fprintf(errOut, "%s%s%s\n", Red, shortener.ErrNotFound.Error(), Reset)
		} else if err != nil {
			fmt.Fprintf(errOut, "%s%s%s\n", Red, ErrSvrError.Error(), Reset)
		} else {
			fmt.Fprintf(out, "%s%s%s\n", Green, result.LongUrl, Reset)
		}

	default:
		fmt.Fprintf(errOut, "%s%s%s%s\n", Red, "unknown command: ", cmd, Reset)
	}
}
