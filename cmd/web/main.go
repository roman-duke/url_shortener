package main

import (
	"context"
	"os"
)

func main() {
	component := hello("john")
	component.Render(context.Background(), os.Stdout)
}
