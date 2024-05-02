package main

import (
	"bufio"
	o "cli/todo/output"
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

var ignoreDirs = []string{".git"}

func UpdateProgress(filename string) {
	fmt.Printf("\rProcessing file: %s", filename)
}

func (c Color) Apply(str string) string {
	return string(c) + str + string(ColorReset)
}

func FindTodos(dir string, output *o.Output) {

	files, err := os.ReadDir(dir)

	if err != nil {
		output.ExitWithError(err.Error(), 2)
	}

	for _, entry := range files {

		filename := entry.Name()
		filepath := fmt.Sprintf("%s/%s", dir, filename)

		if entry.IsDir() {
			if slices.Contains(ignoreDirs, filename) && output.IsTerminal() {
				if output.IsTerminal() {
					output.Printf(ColorRed.Apply("Skipping ignored dir - %s\n"), filename)
				}
				continue
			}
			FindTodos(filepath, output)
		}
		f, err := os.Open(filepath)

		if err != nil {
			output.ExitWithError(err.Error(), 2)
		}
		scanner := bufio.NewScanner(f)

		lineno := 1
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "TODO:") {
				line = strings.Replace(line, "//", "", -1)
				line = strings.Trim(line, " ")
				rawLine := fmt.Sprintf("%d:%s - %s", lineno, filepath, line)

				line = ColorGreen.Apply(rawLine)

				if output.IsTerminal() {
					output.Print(line)
				} else {
					output.Printf("- [ ] %s\n", rawLine)
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

	output := o.CreateOutput()
	if len(os.Args) <= 1 {
		output.ExitWithError("USAGE:  program <folder_name>", 1)
	}
	dir := flag.Args()[0]

	if *exclude != "" {
		ignoreDirs = append(ignoreDirs, *exclude)
	}
	stat, err := os.Stat(dir)

	if err != nil {
		output.ExitWithError(err.Error(), 1)
	}

	if !stat.IsDir() {
		output.ExitWithError(fmt.Sprintf("'%s' is not a directory\n", dir), 1)
	}
	FindTodos(dir, &output)
}
