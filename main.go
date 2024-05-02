package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Color string

const (
	ColorRed   Color = "\033[31m"
	ColorGreen Color = "\033[32m"
	ColorBlue  Color = "\033[34m"
	ColorReset Color = "\033[0m"
)

type Output struct {
	terminal bool
}

func CreateOutput() Output {
	output := Output{
		terminal: false,
	}
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		output.terminal = true
	}
	return output
}
func (o *Output) isTerminal() bool {
	return o.terminal
}

func (o *Output) print(text string) {
	fmt.Println(text)
}

func (o *Output) printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (o *Output) ExitWithError(text string, exitCode int) {
	fmt.Fprintf(os.Stderr, text)
	os.Exit(exitCode)
}

var ignoreDirs = []string{".git"}

func ExitAndPrint(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", e)
		os.Exit(1)
	}
}

func UpdateProgress(filename string) {
	fmt.Printf("\rProcessing file: %s", filename)
}

func (c Color) (str string) string {
	return string(c) + str + string(ColorReset)
}

func FindTodos(dir string, output *Output) {

	files, err := os.ReadDir(dir)

	output.ExitWithError(err.Error())

	for _, entry := range files {

		filename := entry.Name()
		filepath := fmt.Sprintf("%s/%s", dir, filename)

		if entry.IsDir() {
			if slices.Contains(ignoreDirs, filename) && output.isTerminal() {
				if output.isTerminal(){
					output.printf(ColorRed.Apply("Skipping ignored dir - %s\n"), filename)
				}
				continue
			}
			FindTodos(filepath, output)
		}
		f, err := os.Open(filepath)
		output.ExitWithError(err.Error())
		scanner := bufio.NewScanner(f)

		lineno := 1
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "TODO:") {
				line = strings.Replace(line, "//", "", -1)
				line = strings.Trim(line, " ")
				rawLine := fmt.Sprintf("%d:%s - %s", lineno, filepath, line)

				line = ColorGreen.Apply(rawLine)

				if output.isTerminal() {
					output.printf(line)
				} else {
					output.printf("- [ ] %s\n", rawLine)
				}
			}
			lineno++
		}

		f.Close()
	}
}

func main() {

	exclude := flag.String("exclude", "", "Exclude a directory")

	flag.Parse()

	output := CreateOutput()
	if len(os.Args) <= 1 {
		output.ExitWithError("USAGE:  program <folder_name>", 1)
	}
	dir := flag.Args()[0]

	if *exclude != "" {
		ignoreDirs = append(ignoreDirs, *exclude)
	}
	stat, err := os.Stat(dir)

	output.ExitWithError(err.Error(), 1)
	if !stat.IsDir() {
		output.ExitWithError(fmt.Sprintf("'%s' is not a directory\n", dir), 1)
	}
	FindTodos(dir, &output)
}
