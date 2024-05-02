package output

import (
	"fmt"
	"os"
)

type Output struct {
	Terminal bool
}

func CreateOutput() Output {
	output := Output{
		Terminal: false,
	}
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		output.Terminal = true
	}
	return output
}
func (o *Output) IsTerminal() bool {
	return o.Terminal
}

func (o *Output) Print(text string) {
	fmt.Println(text)
}

func (o *Output) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (o *Output) ExitWithError(text string, exitCode int) {
	fmt.Fprintf(os.Stderr, text)
	os.Exit(exitCode)
}
