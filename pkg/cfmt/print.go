package cfmt

import (
	"fmt"
	"io"
	"os"
)

type Color string

const (
	Purple Color = "\033[1;35m%s\033[0m: %s\n"
	Blue   Color = "\033[1;34m%s\033[0m: %s\n"
	Green  Color = "\033[1;32m%s\033[0m: %s\n"
	Red    Color = "\033[1;31m%s\033[0m: %s\n"
	Orange Color = "\033[1;33m%s\033[0m: %s\n"
)

func FPrint(r io.Writer, color Color, pref, s string) (int, error) {
	return fmt.Fprintf(r, string(color), pref, s)
}

func Print(color Color, pref, s string) (int, error) {
	return FPrint(os.Stdout, color, pref, s)
}
